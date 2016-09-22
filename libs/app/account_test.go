package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
)

func Test_App(t *testing.T) {
	Convey("ACCOUNT COMMANDS", t, func() {
		app := CreateAppStub()
		u := app.Users.(*users.UsersMock)
		t := app.Settings.(*settings.SettingsMock)
		w := app.TemplateWriter.(*gui.TemplateWriterMock)
		d := app.Dialogues.(*gui.DialogueMock)

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
			Convey(`Should call api signin and enter form details`, func() {
				d.FieldResponse = &gui.DialogueResponse{Value:"richard"}
				app.App.Run([]string{"[app]", "signin"})
				So(d.FieldCallHistory[0], should.Resemble, []interface{}{"account:usernamefield", nil})
				So(d.FieldCallHistory[1], should.Resemble, []interface{}{"account:passwordfield", nil})
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

		Convey(`cd`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("cd")
				So(p, should.NotBeNil)
			})
			Convey(`Should change directory`, func() {
				app.App.Run([]string{"[app]", "cd", "dillbuck"})
				So(t.ChangeDirectoryCalledWith, should.Equal, "dillbuck")
				So(w.RenderCalledWith, should.Resemble, []interface{}{"account:cd", map[string]string{"username" :"dillbuck"}})
			})
		})
	})
}


