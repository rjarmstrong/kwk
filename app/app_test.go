package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/api"
)

func Test_System(t *testing.T) {
	Convey("Manage settings", t, func() {
		a := &api.ApiMock{}
		app := NewKwkApp(a)

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
		})

	})
}


