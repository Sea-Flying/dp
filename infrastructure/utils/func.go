package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetExecPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//merge two slice of unique elements into one
func MergeStringSlice(a []string, b []string) (ret []string) {
	ret = a
	for _, i := range b {
		var notExisted = true
		for _, j := range a {
			if j == i {
				notExisted = false
				break
			}
		}
		if notExisted {
			ret = append(ret, i)
		}
	}
	return
}

func IsInSlice(s []string, e string) bool {
	for _, i := range s {
		if i == e {
			return true
		}
	}
	return false
}
