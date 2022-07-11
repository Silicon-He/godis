package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsExist(err)
}

func checkPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func isNotExistMkDir(src string) error {
	if !checkExist(src) {
		if err := os.MkdirAll(src, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func mustOpen(fileName, dir string) (f *os.File, err error) {
	if checkPermission(dir) {
		return nil, fmt.Errorf("permission denined dir:%s", dir)
	}
	if err = isNotExistMkDir(dir); err != nil {
		return nil, fmt.Errorf("error during make dir:%s", dir)
	}
	f, err = os.OpenFile(filepath.Join(dir, fileName), os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("fail to open file,err :%s", err)
	}
	return f, nil
}
