package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
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
	cmd.Dir = "/root/busybox"

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

// pivotRoot 改变挂载的位置
func pivotRoot(rootpwd string) (err error) {
	var pivotDir string // 存储 old_root 的路径

	// 使当前rootPwd的老rootPwd和新rootPwd不在同一个文件系统下，重新mount一次rootPwd
	// bind mount 将相同的内容换了一个挂载点
	if err = syscall.Mount(rootpwd, rootpwd, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	// 创建 rootPwd/.pivot_root 存储旧 old_root
	pivotDir = filepath.Join(rootpwd, ".pivot_root")
	if err = os.Mkdir(pivotDir, 0777); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("os.Mkdir( %s ) error %v", pivotDir, err)
		}
	}

	// 更换挂载位置
	if err = syscall.PivotRoot(rootpwd /* newroot */, pivotDir /* putold */); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// 切换工作目录
	if err = syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	// unmount old_root
	pivotDir = filepath.Join("/", ".pivot_root")
	if err = syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	// 删除无用的目录
	if err = os.Remove(pivotDir); err != nil {
		return fmt.Errorf("Remove dir %s error, %v", pivotDir, err)
	}

	log.Printf("mount child process on %s\n", rootpwd)

	return
}
