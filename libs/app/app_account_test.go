package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/gui"
)

func Test_App(t *testing.T) {
	Convey("ACCOUNT COMMANDS", t, func() {
		a := &api.ApiMock{}
		w := &gui.Writer{}
		app := NewKwkApp(a, nil, w, nil)

		Convey(`Profile`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("profile")
				So(p, should.NotBeNil)
				p2 := app.App.Command("me")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call api print profile`, func() {
				app.App.Run([]string{"[app]", "profile"})
				So(a.PrintProfileCalled, should.BeTrue)
			})
		})

		Convey(`Signin`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("signin")
				So(p, should.NotBeNil)
				p2 := app.App.Command("login")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call api signin`, func() {
				app.App.Run([]string{"[app]", "signin", "richard", "password"})
				So(a.LoginCalledWith[0], should.Equal, "richard")
				So(a.LoginCalledWith[1], should.Equal, "password")
			})
		})

		Convey(`Signup`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("signup")
				So(p, should.NotBeNil)
				p2 := app.App.Command("register")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call api signup`, func() {
				app.App.Run([]string{"[app]", "signup", "richard@kwk.co", "richard", "password"})
				So(a.SignupCalledWith[0], should.Equal, "richard@kwk.co")
				So(a.SignupCalledWith[1], should.Equal, "richard")
				So(a.SignupCalledWith[2], should.Equal, "password")
			})
		})

		Convey(`Signout`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("signout")
				So(p, should.NotBeNil)
				p2 := app.App.Command("logout")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call api signout`, func() {
				app.App.Run([]string{"[app]", "signout"})
				So(a.SignoutCalled, should.BeTrue)
			})
		})
	})
}


