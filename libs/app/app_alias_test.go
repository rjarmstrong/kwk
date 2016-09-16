package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/openers"
)

func Test_Alias(t *testing.T) {
	Convey("ALIAS COMMANDS", t, func() {
		a := &api.ApiMock{}
		w := &gui.WriterMock{}
		s := &system.SystemMock{}
		o := &openers.OpenerMock{}
		app := NewKwkApp(a, s, w, o)

		Convey(`Command not found`, func() {
			Convey(`Should call get and open if found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []api.KwkLink{
					{FullKey:fullKey},
				}
				app.App.Run([]string{"[app]", fullKey, "arg1", "arg2"})
				So(a.GetCalledWith, should.Equal, fullKey)
				So(o.OpenCalledWith, should.Resemble, []interface{}{&a.ReturnItemsForGet[0], []string{"arg1", "arg2"}})
			})
			Convey(`Should call get and prompt if multiple found`, func() {
			})
			Convey(`Should call get print if not found`, func() {
			})
		})

		Convey(`New`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("new")
				So(p, should.NotBeNil)
				p2 := app.App.Command("create")
				p3 := app.App.Command("save")
				So(p2.Name, should.Equal, p.Name)
				So(p3.Name, should.Equal, p.Name)
			})
			Convey(`Should call create, save to clip board and print with template WITH a fullKey`, func() {
				fullKey := "hola.sh"
				app.App.Run([]string{"[app]", "new", "echo hola", fullKey})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", fullKey})
				So(s.CopyToClipboardCalledWith, should.Equal, fullKey)
				So(w.PrintCalledWith, should.Resemble, []interface{}{"new", &api.KwkLink{FullKey:fullKey}})
			})
			Convey(`Should call create, save to clip board and print with template WITHOUT a fullKey`, func() {
				app.App.Run([]string{"[app]", "new", "echo hola"})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", ""})
				mockKey := "x5hi23"
				So(s.CopyToClipboardCalledWith, should.Equal, mockKey)
				So(w.PrintCalledWith, should.Resemble, []interface{}{"new", &api.KwkLink{FullKey:mockKey}})
			})
		})

		Convey(`Inspect`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("inspect")
				So(p, should.NotBeNil)
				p2 := app.App.Command("i")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call get and print with template`, func() {
				app.App.Run([]string{"[app]", "inspect", "arrows.js"})
				So(a.GetCalledWith, should.Equal, "arrows.js")
				So(w.PrintCalledWith, should.Resemble, []interface{}{"inspect", &api.KwkLinkList{}})
			})
		})

		Convey(`Cat`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("cat")
				So(p, should.NotBeNil)
				p2 := app.App.Command("raw")
				p3 := app.App.Command("read")
				p4 := app.App.Command("print")
				p5 := app.App.Command("get")
				So(p2.Name, should.Equal, p.Name)
				So(p3.Name, should.Equal, p.Name)
				So(p4.Name, should.Equal, p.Name)
				So(p5.Name, should.Equal, p.Name)
			})
			Convey(`Should call get and print with template`, func() {
				app.App.Run([]string{"[app]", "cat", "arrows.js"})
				So(a.GetCalledWith, should.Equal, "arrows.js")
				So(w.PrintCalledWith, should.Resemble, []interface{}{"cat", &api.KwkLinkList{}})
			})
		})
	})
}


