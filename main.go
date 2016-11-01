package main

import (
	"os"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/app"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"bufio"
)

func main() {
	os.Setenv("version", "v0.0.3")

	conn := rpc.Conn("aus1.kwk.co:80");
	defer conn.Close()

	s := system.New()
	t := settings.New(s, "settings")
	h := rpc.NewHeaders(t)
	u := users.New(conn, t, h)
	a := aliases.New(conn, t, h)
	o := openers.New(s, a)
	w := gui.NewTemplateWriter(os.Stdout)
	r := bufio.NewReader(os.Stdin)
	d := gui.NewDialogues(w, r)

	kwkApp := app.NewKwkApp(a, s, t, o, u, d, w)
	kwkApp.App.Run(os.Args)
}
