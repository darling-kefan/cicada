package cicada

import (
	"os"
	"os/exec"
	"path/filepath"
)

// 工具函数集

func runPath() string {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return path
}


