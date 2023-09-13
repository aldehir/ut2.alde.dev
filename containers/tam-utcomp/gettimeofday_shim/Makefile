libgettimeofday_shim.so: gettimeofday_shim.o
	gcc -m32 -shared -o $@ $^

%.o : %.c
	gcc -c -m32 -fPIC $< -o $@
