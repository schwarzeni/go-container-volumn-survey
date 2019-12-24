package main

import (
	"fmt"
	"go-volumn/aufs"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	rootURL := "/root/workplace_go/go-volumn-dev"
	mntURL := "/root/workplace_go/go-volumn-dev/mnt"
	defer aufs.DeleteWorkSpace(rootURL, mntURL)
	if os.Args[0] == "/proc/self/exe" { // child process
		childProcess()
		return
	}

	var (
		cmd *exec.Cmd
		err error
	)
	cmd = exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Dir = mntURL
	if err = aufs.NewWorkSpace(rootURL, mntURL); err != nil {
		log.Fatal(err)
	}

	if err = cmd.Start(); err != nil {
		log.Fatalf("cmd.Start() failed: %v", err)
	}
	cmd.Wait()
}

func childProcess() {
	var (
		pwd string
		err error
	)
	if pwd, err = os.Getwd(); err != nil {
		fmt.Fprintf(os.Stderr, "Get current location error %v", err)
		return
	}

	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	if err = pivotRoot(pwd); err != nil {
		fmt.Fprintf(os.Stderr, "pivotRoot( %s ) error %v", pwd, err)
		return
	}
	if err = syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, ""); err != nil {
		fmt.Fprintf(os.Stderr, "mount proc error %v", err)
		return
	}
	if err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=0755"); err != nil {
		fmt.Fprintf(os.Stderr, "mount tmpfs error %v", err)
		return
	}
	if err := syscall.Exec("/bin/sh", []string{"sh"}, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "exec error %v", err)
		return
	}
}
