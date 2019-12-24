# 《自己动手写Docker》第四章相关技术调研

---

使用 busybox 镜像作为样例，在克隆出的子进程中运行 `/bin/sh` 程序，相关地方使用了硬编码

- pivot_root
- mount aufs
- volume

环境为：

- Linux Ubuntu 16.04
- 内核版本 4.10.0-28-generic
- Go 1.12.4

---

## pivot_root

Go语言为我们封装的 `PivotRoot` 函数原型就很好说明了其作用：

```go
func PivotRoot(newroot string, putold string) (err error)
```

执行这个函数之前，还需要执行如下函数，否则会报错。[参考1](https://github.com/xianlubird/mydocker/issues/13#issuecomment-450898307)，[参考2](https://github.com/xianlubird/mydocker/issues/41#issuecomment-478799767)

```go
syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
```

将当前整个系统系统切换到一下新的目录下(newroot)，移除就文件系统至另一个目录(putold)。那么请注意，在执行这个之前必须对该进程的 Mount Namespace 进行隔离，否则的话整个系统就会挂掉了，所以需要如下代码来保证命名空间的隔离

```go
	cmd = exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID}
```

其它的流程就像 [pivot_root.go](./pivot_root.go) 中所注释的那样，该操作由克隆出的子进程执行，执行完毕后使用 `syscall.Exec` 将其替换成真正需要执行的程序

---

## mount aufs

主要实现在 [aufs/new_workspace.go](./aufs/new_workspace.go) 下，排除冗余的逻辑，核心就是执行一条 `mount` 命令，示例如下：

```sh
mount -t aufs -o dirs=./container-layer:./image-layer1:./image-layer2 none ./mnt
# none: specifies we don’t have any device associated with it, since we are going to mount two directories
```

`container-layer` 中为可写的，之后的文件夹中的内容都是不可写的，使用冒号分割，最终的目标为 `mnt` 文件夹

当然，要有始有终，结束后需要清除相关的可写层以及挂载的文件，主要实现在 [aufs/del_workspace.go](./aufs/del_workspace.go) 中

---

## volume

这个讲的是自定义一些目录挂载至容器中，也是用的和上一部分相似的技术，主要实现在 [aufs/mnt.go](./aufs/mnt.go) 中

```go
	cmd := exec.Command("mount", "-t", "aufs", "-o",
		fmt.Sprintf("dirs=%s", srcURL), "none", path.Join(mntURL, dstURL))
```
