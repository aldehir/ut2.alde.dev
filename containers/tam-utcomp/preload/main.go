package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ucc = "ucc-bin"

func init() {
	envUCC, ok := os.LookupEnv("UCC")
	if ok {
		ucc = envUCC
	}
}

func main() {
	Must1(os.Chdir("System"))

	cmd := exec.Command(ucc, "server", "-nohomedir")
	stdoutPipe := Must2(cmd.StdoutPipe())
	cmd.Stderr = os.Stderr

	done := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)

			if strings.Contains(line, "UdpServerQuery") {
				done <- struct{}{}
			}
		}
	}()

	Must1(cmd.Start())

	timer := time.NewTimer(10 * time.Minute)

	var success bool

loop:
	for {
		select {
		case <-done:
			success = true
			break loop
		case <-timer.C:
			break loop
		}
	}

	close(done)

	cmd.Process.Kill()
	cmd.Wait()

	if success {
		fmt.Println("Preload complete.")
	} else {
		fmt.Println("Preload timed out.")
		os.Exit(1)
	}

	os.Exit(0)
}

func Must1(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
