package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/app"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"bytes"
	"github.com/smartystreets/assertions/should"
)

func Test_App(t *testing.T) {
	Convey("ACCOUNT COMMANDS", t, func() {
		conn := rpc.Conn("127.0.0.1:7777");
		defer conn.Close()
		s := system.New()
		t := settings.New(s, "settings")
		h := rpc.NewHeaders(t)
		u := users.New(conn, t, h)
		a := aliases.New(conn, t, h)
		o := openers.New(s, a)
		b := &bytes.Buffer{}
		w := gui.NewTemplateWriter(b)
		d := gui.NewDialogues(w)

		kwk := app.NewKwkApp(a, s, t, o, u, d, w)
		Convey(`Profile`, func() {
			Convey(`Should print not logged in`, func() {
				kwk.App.Run([]string{"[app]", "profile"})
				So(b.String(), should.Equal, "You are not logged in please log in: kwk login <username> <password>\n")
				b.Reset()
			})
			Convey(`Should print profile`, func() {
				kwk.App.Run([]string{"[app]", "signin", "richard", "D1llbuck"})
				b.Reset()
				kwk.App.Run([]string{"[app]", "profile"})
				So(b.String(), should.Equal, "You are: richard!")
				kwk.App.Run([]string{"[app]", "signout"})
			})
		})
	})
}
