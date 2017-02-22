package models

import "fmt"

type ApiInfo struct {
	Version         string
	Build           string
	Revision        string
	ClientSupported string
	Dc              string
}

func (a *ApiInfo) String() string {
	return fmt.Sprintf("%s+%s", a.Version, a.Build)
}

var Api ApiInfo = ApiInfo{Version:"v-.-.-", Build:"0"}