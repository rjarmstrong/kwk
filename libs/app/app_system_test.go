package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/gui"
)

func Test_System(t *testing.T) {
	Convey("SYSTEM COMMANDS", t, func() {
		s := &system.SystemMock{}
		w := &gui.WriterMock{}
		app := NewKwkApp(nil, s, w)

		Convey(`Upgrade`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("upgrade")
				So(p, should.NotBeNil)
			})
			Convey(`Should call upgrade`, func() {
				app.App.Run([]string{"[app]", "upgrade"})
				So(s.UpgradeCalled, should.BeTrue)
			})
		})

		Convey(`Version`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("version")
				So(p, should.NotBeNil)
			})
			Convey(`Should get version and call writer print`, func() {
				app.App.Run([]string{"[app]", "version"})
				So(s.VersionCalled, should.BeTrue)
				So(w.PrintCalledWith, should.Equal, "0.0.1")
			})
		})

		Convey(`CD`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("cd")
				So(p, should.NotBeNil)
			})
			Convey(`Should change directory`, func() {
				app.App.Run([]string{"[app]", "cd", "dillbuck"})
				So(s.ChangeDirectoryCalledWith, should.Equal, "dillbuck")
				So(w.PrintCalledWith, should.Equal, "Changed to dillbuck")
			})
		})
	})
}


