package cache

import (
	"os/user"
	"runtime"
	"os"
	"fmt"
)

const (
	StandardFilePermission = 0700
	OS_DARWIN = `darwin`
	OS_LINUX = `linux`
	OS_WINDOWS = `windows`
)

var cachePath string

func Path() string {
	if cachePath != "" {
		return cachePath
	}
	var p string
	u, _ := user.Current()
	p = fmt.Sprintf("%s/.kwk", u.HomeDir)
	if runtime.GOOS == OS_WINDOWS {
		p = "%LocalAppData%\\kwk"
	} else if runtime.GOOS == OS_DARWIN {
		if u.Username == "root" {
			p = "/var/root/.kwk"
		}
	}
	if err := os.MkdirAll(p, StandardFilePermission); err != nil {
		if os.IsExist(err) {
			return p
		}
		fmt.Println("kwk might not be able to run on your system please copy error and report to: https://github.com/kwk-cli/cli-issues", err)
		return ""
	}
	cachePath = p
	return cachePath
}
