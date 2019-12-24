package aufs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

// NewWorkSpace 创建新的容器层
func NewWorkSpace(rootURL string, mntURL string, volume string) (err error) {
	if err = createReadOnlyLayer(rootURL); err != nil {
		goto ERR
	}
	if err = createWriteLayer(rootURL); err != nil {
		goto ERR
	}
	if err = createMountPoint(rootURL, mntURL); err != nil {
		goto ERR
	}
	if len(volume) != 0 {
		var volumeURLs []string
		if volumeURLs, err = volumeURLExtract(volume); err != nil {
			goto ERR
		}
		if err = mountVolume(rootURL, mntURL, volumeURLs); err != nil {
			goto ERR
		}
		log.Printf("%s mount on %s in container", volumeURLs[0], volumeURLs[1])
	}
	return
ERR:
	deleteWriteLayer(rootURL)
	deleteMountPoint(rootURL, mntURL)
	return
}

// createReadOnlyLayer 解压镜像，创建容器只读层
func createReadOnlyLayer(rootURL string) (err error) {
	var exist bool
	busyboxURL := path.Join(rootURL, "busybox")
	busyboxTarURL := path.Join(rootURL, "busybox.rar")
	exist, err = pathExists(busyboxURL)
	if err != nil {
		return fmt.Errorf("Failed to judge whether dir %s exists. %v", busyboxURL, err)
	}
	if exist == false {
		if err = os.Mkdir(busyboxURL, 0777); err != nil {
			return fmt.Errorf("Midir dir %s error. %v", busyboxURL, err)
		}
		if _, err = exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			return fmt.Errorf("unTar dir %s error %v", busyboxTarURL, err)
		}
	}
	return
}

// createWriteLayer 创建容器可写层
func createWriteLayer(rootURL string) (err error) {
	writeURL := path.Join(rootURL, "writeLayer/")
	if err = os.Mkdir(writeURL, 0777); err != nil {
		return fmt.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
	return
}

// createMountPoint 将只读层和可写层都mount到一处
func createMountPoint(rootURL string, mntURL string) (err error) {
	if err = os.Mkdir(mntURL, 0777); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("Mkdir dir %s error. %v", mntURL, err)
		}
		err = nil
	}
	dirs := fmt.Sprintf("dirs=%s:%s", path.Join(rootURL, "writeLayer"), path.Join(rootURL, "busybox"))
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd %s %v error %v", cmd.Path, cmd.Args, err)
	}
	return
}
