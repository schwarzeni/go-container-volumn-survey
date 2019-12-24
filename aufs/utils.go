package aufs

import (
	"fmt"
	"os"
	"strings"
)

// pathExists 判断路径是否存在
func pathExists(path string) (exist bool, err error) {
	if _, err = os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// volumeURLExtract 解析参数 src:dst
func volumeURLExtract(volume string) ([]string, error) {
	var (
		volumeURLs []string
	)
	volumeURLs = strings.Split(volume, ":")
	if len(volumeURLs) != 2 || len(volumeURLs[0]) == 0 || len(volumeURLs[1]) == 0 {
		return nil, fmt.Errorf("volumeURLExtract failed: %s", volume)
	}
	return volumeURLs, nil
}
