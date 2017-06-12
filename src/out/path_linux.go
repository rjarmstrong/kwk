package out

import (
	"fmt"
	"os/user"
)

func getPath() string {
	u, _ := user.Current()
	return fmt.Sprintf("%s/.kwk", u.HomeDir)
}
