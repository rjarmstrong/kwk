package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_System(t *testing.T) {
	Convey("SYSTEM COMMANDS", t, func() {
		app := CreateAppStub()
		s := app.System.(*system.SystemMock)
		w := app.TemplateWriter.(*gui.TemplateWriterMock)

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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"system:version", map[string]string{
					"version": "0.0.1",
				}})
			})
		})
	})
}
