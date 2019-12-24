package aufs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

// DeleteWorkSpace 删除容器层
func DeleteWorkSpace(rootURL string, mntURL string, volume string) (err error) {
	if len(volume) != 0 {
		var volumeURLs []string
		if volumeURLs, err = volumeURLExtract(volume); err != nil {
			goto ERR
		}
		if err = umountVolume(rootURL, mntURL, volumeURLs); err != nil {
			goto ERR
		}
	}
	if err = deleteMountPoint(rootURL, mntURL); err != nil {
		goto ERR
	}
	if err = deleteWriteLayer(rootURL); err != nil {
		goto ERR
	}
	return
ERR:
	// TODO: handle error here
	return
}

// deleteMountPoint 删除挂载点
func deleteMountPoint(rootURL string, mntURL string) (err error) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("unmount %s error %v", mntURL, err)
	}
	if err = os.RemoveAll(mntURL); err != nil {
		return fmt.Errorf("removeAll %s error %v", mntURL, err)
	}
	return
}

// deleteWriteLayer 删除可写层
func deleteWriteLayer(rootURL string) (err error) {
	writeLayer := path.Join(rootURL, "writeLayer")
	if err = os.RemoveAll(writeLayer); err != nil {
		return fmt.Errorf("removeAll %s error %v", writeLayer, err)
	}
	return
}
