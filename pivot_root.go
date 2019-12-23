package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

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
