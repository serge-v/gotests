/* ivan(a.t)mysqlab.net */

package main

import (
        "syscall"
        "os"
        "log"
        "fmt"
)

func daemon(nochdir, noclose int) int {
        var ret, ret2 uintptr
        var errno syscall.Errno
	var err error

        // already a daemon
        if syscall.Getppid() == 1 {
		log.Println("already a daemon")
                return 0
        }

        // fork off the parent process
        ret, ret2, errno = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
        if errno != 0 {
		log.Println("fork: ", errno)
                return -1
        }

	pid := syscall.Getpid()

	fmt.Println("fork: ", ret, ret2, "pid:", pid)

	// failure
        if ret2 < 0 {
                os.Exit(-1)
        }

        // if we got a good PID, then we call exit the parent process.
        if ret > 0 {
		fmt.Println("good pid: ", ret, "pid:", pid)
                os.Exit(0)
        }

        /* Change the file mode mask */
        _ = syscall.Umask(0)

        // create a new SID for the child process
        s_ret, err := syscall.Setsid()
	fmt.Println("set side: ",s_ret, err, "pid:", pid)
        if err != nil {
                log.Printf("Error: syscall.Setsid error: %s", err)
        }
        if s_ret < 0 {
                return -1
        }

        if nochdir == 0 {
                os.Chdir("/")
        }

        if noclose == 0 {
                f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
                if e == nil {
                        fd := f.Fd()
                        syscall.Dup2(int(fd), int(os.Stdin.Fd()))
                        syscall.Dup2(int(fd), int(os.Stdout.Fd()))
                        syscall.Dup2(int(fd), int(os.Stderr.Fd()))
                }
        }

        return 0
}

// usage example: daemon(0, 0)
