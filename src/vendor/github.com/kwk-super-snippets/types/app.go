package types

import "fmt"

type AppInfo struct {
	Version string
	Build   string
	Time    int64
	Notes   string
}

func (c *AppInfo) String() string {
	return fmt.Sprintf("%s+%s", c.Version, c.Build)
}
