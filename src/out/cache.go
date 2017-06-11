package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"os"
	"os/user"
	"runtime"
)

const (
	StandardFilePermission = 0700
)

var kwkPath string

// TODO: CHANGE TO STATIC PATH RESOLUTION
func KwkPath() string {
	if kwkPath != "" {
		return kwkPath
	}
	var p string
	u, _ := user.Current()
	p = fmt.Sprintf("%s/.kwk", u.HomeDir)
	if runtime.GOOS == types.OsWindows {
		p = "%LocalAppData%\\kwk"
	} else if runtime.GOOS == types.OsDarwin {
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
	kwkPath = p
	return kwkPath
}
