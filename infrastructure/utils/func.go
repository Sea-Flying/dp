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
