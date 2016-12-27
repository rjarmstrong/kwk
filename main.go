package main

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/app"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"bufio"
	"os"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/dlg"
)

func main() {
	os.Setenv("version", "v0.0.3")
	host := os.Getenv("KWK_HOST")
	conn := rpc.Conn(host)
	defer conn.Close()

	s := system.New()
	t := settings.New(s, "settings")
	h := rpc.NewHeaders(t)
	u := users.New(conn, t, h)
	a := aliases.New(conn, t, h)
	w := tmpl.NewWriter(os.Stdout)
	o := openers.New(s, a, w)
	r := bufio.NewReader(os.Stdin)
	d := dlg.New(w, r)
	ch := search.New(conn, t, h)

	kwkApp := app.NewKwkApp(a, s, t, o, u, d, w, ch)
	kwkApp.App.Run(os.Args)
}
