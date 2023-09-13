#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <dlfcn.h>

void *libc;
int (*__gettimeofday)(struct timeval * restrict tv, void * restrict tz);

void __attribute__((constructor)) gettimeofday_init(void) {
  puts("Initializing gettimeofday");
}

void __attribute__((destructor)) gettimeofday_deinit(void) {
  dlclose(libc);
}

static void load_original() {
  libc = dlopen("libc.so.6", RTLD_LAZY);
  if (!libc) {
    fprintf(stderr, "%s\n", dlerror());
    exit(EXIT_FAILURE);
  }

  __gettimeofday = dlsym(libc, "gettimeofday");
  if (!__gettimeofday) {
    fprintf(stderr, "%s\n", dlerror());
    exit(EXIT_FAILURE);
  }
}

int gettimeofday(struct timeval * restrict tv, void * restrict tz) {
  if (tz != NULL) {
    return __gettimeofday(tv, tz);
  }

  struct timespec ts;
  if (clock_gettime(CLOCK_MONOTONIC, &ts) == 0) {
    tv->tv_sec = ts.tv_sec;
    tv->tv_usec = ts.tv_nsec / 1000;
    return 0;
  }

  return -1;
}
