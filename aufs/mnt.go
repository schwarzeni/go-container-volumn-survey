package aufs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func mountVolumn(rootURL string, mntURL string, volumnURLs []string) (err error) {
	var (
		srcURL     = volumnURLs[0]
		dstURL     = volumnURLs[1]
		fulldstURL = path.Join(mntURL, dstURL)
	)

	if err := os.Mkdir(srcURL, 0777); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("os mkdir srcURL %s error, %v", srcURL, err)
		}
	}

	if err := os.Mkdir(fulldstURL, 0777); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("os mkdir fulldstURL %s error, %v", srcURL, err)
		}
	}

	cmd := exec.Command("mount", "-t", "aufs", "-o",
		fmt.Sprintf("dirs=%s", srcURL), "none", fulldstURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Mount volume failed. %v", err)
	}

	return
}
