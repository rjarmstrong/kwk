package main

import (
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
	"os"
	"github.com/kwk-links/kwk-cli/libs/services/openers"
	"github.com/kwk-links/kwk-cli/libs/app"
	"github.com/kwk-links/kwk-cli/libs/rpc"
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"github.com/kwk-links/kwk-cli/libs/services/system"
)

func main() {
	os.Setenv("version", "v0.0.1")

	conn := rpc.Conn("127.0.0.1:7777");
	defer conn.Close()

	s := system.New()
	t := settings.New(s, "settings")

	u := users.New(conn, t)
	a := aliases.New(conn, t)
	o := openers.New(s, a)

	kwkApp := app.NewKwkApp(a, s, t, o, u)
	kwkApp.App.Run(os.Args)
}
