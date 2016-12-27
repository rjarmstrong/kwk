package main

import (
	"bitbucket.com/sharingmachine/kwkcli/app"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/openers"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/settings"
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/users"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bufio"
	"os"
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
	a := snippets.New(conn, t, h)
	w := tmpl.NewWriter(os.Stdout)
	o := openers.New(s, a, w)
	r := bufio.NewReader(os.Stdin)
	d := dlg.New(w, r)
	ch := search.New(conn, t, h)

	kwkApp := app.New(a, s, t, o, u, d, w, ch)
	kwkApp.App.Run(os.Args)
}
