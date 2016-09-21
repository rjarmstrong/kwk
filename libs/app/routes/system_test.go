package routes

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/app"
)

func Test_System(t *testing.T) {
	Convey("SYSTEM COMMANDS", t, func() {
		s := &system.SystemMock{}
		sett := &settings.SettingsMock{}
		appl := app.NewKwkApp(nil, s, sett, nil, nil)

		Convey(`Upgrade`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("upgrade")
				So(p, should.NotBeNil)
			})
			Convey(`Should call upgrade`, func() {
				appl.App.Run([]string{"[app]", "upgrade"})
				So(s.UpgradeCalled, should.BeTrue)
			})
		})

		Convey(`Version`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("version")
				So(p, should.NotBeNil)
			})
			Convey(`Should get version and call writer print`, func() {
				appl.App.Run([]string{"[app]", "version"})
				So(s.VersionCalled, should.BeTrue)
				//So(i.LastRespondCalledWith, should.Resemble, []interface{}{"version", "0.0.1"})
			})
		})

		Convey(`CD`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("cd")
				So(p, should.NotBeNil)
			})
			Convey(`Should change directory`, func() {
				appl.App.Run([]string{"[app]", "cd", "dillbuck"})
				So(sett.ChangeDirectoryCalledWith, should.Equal, "dillbuck")
				//So(i.LastRespondCalledWith, should.Resemble, []interface{}{"cd", "dillbuck"})
			})
		})
	})
}


