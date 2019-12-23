package aufs

import "os"

func pathExists(path string) (exist bool, err error) {
	if _, err = os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
