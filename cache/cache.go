package cache

import (
	"os/user"
	"runtime"
	"fmt"
	"os"
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
	if runtime.GOOS == OS_WINDOWS {
		// use AppDir instead
		p = "%LocalAppData%\\kwk"
	} else if runtime.GOOS == OS_LINUX {
		p = fmt.Sprintf("/%s/.kwk", u.Username)
	} else if runtime.GOOS == OS_DARWIN {
		if u.Username == "root" {
			p = "/var/root/.kwk"
		} else {
			p = fmt.Sprintf("/Users/%s/.kwk", u.Username)
		}
	} else {
		// TODO: Write friendly
		panic("OS not supported.")
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
