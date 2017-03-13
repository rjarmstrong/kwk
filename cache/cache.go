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
	u, err := user.Current()
	if err != nil {
		// TODO: Write friendly
		panic(err)
	}
	p = fmt.Sprintf("%s/.kwk", u.HomeDir)
	if runtime.GOOS == OS_WINDOWS {
		// use AppDir instead
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
		panic(err)
	}
	cachePath = p
	return cachePath
}
