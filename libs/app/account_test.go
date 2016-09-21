package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/users"
)

func Test_App(t *testing.T) {
	Convey("ACCOUNT COMMANDS", t, func() {
		u := &users.UsersMock{}
		sett := &settings.SettingsMock{}
		app := NewKwkApp(nil, nil, sett, nil, u)

		Convey(`Profile`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("profile")
				So(p, should.NotBeNil)
				p2 := app.App.Command("me")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should get profile and respond with template`, func() {
				app.App.Run([]string{"[app]", "profile"})
				So(u.GetCalled, should.BeTrue)
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
				So(u.LoginCalledWith[0], should.Equal, "richard")
				So(u.LoginCalledWith[1], should.Equal, "password")
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
				So(u.SignupCalledWith[0], should.Equal, "richard@kwk.co")
				So(u.SignupCalledWith[1], should.Equal, "richard")
				So(u.SignupCalledWith[2], should.Equal, "password")
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
				So(u.SignoutCalled, should.BeTrue)
			})
		})
	})
}


