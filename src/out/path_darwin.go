package out

import (
	"fmt"
	"os/user"
)

func getPath() string {
	u, _ := user.Current()
	if u.Username == "root" {
		return "/var/root/.kwk"
	}
	return fmt.Sprintf("%s/.kwk", u.HomeDir)
}
