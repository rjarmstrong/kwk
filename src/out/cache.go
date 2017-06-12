package out

import (
	"fmt"
	"os"
)

const (
	// StandardFilePermission is the permission all files are stored in ~/.kwk by default.
	StandardFilePermission = 0700
)

var kwkPath string

// KwkPath is the root directory of the local kwk cache
func KwkPath() string {
	if kwkPath != "" {
		return kwkPath
	}
	p := getPath()
	if err := os.MkdirAll(p, StandardFilePermission); err != nil {
		if os.IsExist(err) {
			return p
		}
		fmt.Println("kwk might not be able to run on your system please "+
			"copy error and report at: https://github.com/rjarmstrong/kwk", err)
		return ""
	}
	kwkPath = p
	return kwkPath
}
