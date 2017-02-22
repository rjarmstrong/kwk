package models

import (
	"fmt"
	"time"
)

type ClientInfo struct {
	Version string
	Build   string
	Time    int64
	Notes   string
	Api ApiInfo
}

var Client = ClientInfo{Api: Api}

func (c *ClientInfo) String() string {
	return fmt.Sprintf("%s+%s", c.Version, c.Build)
}

func ClientIsNew(t int64) bool {
	if t == 0 {
		return false
	}
	return t > (time.Now().Unix()-60)
}