package app

import (
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_System(t *testing.T) {
	Convey("SYSTEM CLI", t, func() {
		app := CreateAppStub()
		s := app.System.(*system.MockSystem)
		w := app.TemplateWriter.(*tmpl.MockWriter)

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
