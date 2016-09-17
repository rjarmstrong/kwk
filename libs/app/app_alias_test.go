package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/openers"
	"github.com/iris-contrib/errors"
)

func Test_Alias(t *testing.T) {
	Convey("ALIAS COMMANDS", t, func() {
		a := &api.ApiMock{}
		i := &gui.InteractionMock{}
		s := &system.SystemMock{}
		o := &openers.OpenerMock{}
		app := NewKwkApp(a, s, i, o)

		Convey(`Command not found`, func() {
			Convey(`Should call get and open if found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []api.Alias{
					{FullKey:fullKey},
				}
				app.App.Run([]string{"[app]", fullKey, "covert", "arg2"})
				So(a.GetCalledWith, should.Equal, fullKey)
				So(o.OpenCalledWith, should.Resemble, []interface{}{&a.ReturnItemsForGet[0], []string{"covert", "arg2"}})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should call 'get' and prompt if multiple found`, func() {
				fullKey1 := "hola.sh"
				fullKey2 := "hola.js"
				a.ReturnItemsForGet = []api.Alias{
					{FullKey:fullKey1},
					{FullKey:fullKey2},
				}
				app.App.Run([]string{"[app]", "hola", "arg1", "arg2"})
				So(a.GetCalledWith, should.Equal, "hola")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"chooseBetweenKeys", a.ReturnItemsForGet})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should respond if not found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []api.Alias{}
				app.App.Run([]string{"[app]", fullKey, "arg1", "arg2"})
				So(a.GetCalledWith, should.Equal, fullKey)
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"notfound", fullKey})
				a.ReturnItemsForGet = nil
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
			Convey(`Should call create, save to clip board and respond with template WITH a fullKey`, func() {
				fullKey := "hola.sh"
				app.App.Run([]string{"[app]", "new", "echo hola", fullKey})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", fullKey})
				So(s.CopyToClipboardCalledWith, should.Equal, fullKey)
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"new", &api.Alias{FullKey:fullKey}})
			})
			Convey(`Should call create, save to clip board and respond with template WITHOUT a fullKey`, func() {
				app.App.Run([]string{"[app]", "new", "echo hola"})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", ""})
				mockKey := "x5hi23"
				So(s.CopyToClipboardCalledWith, should.Equal, mockKey)
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"new", &api.Alias{FullKey:mockKey}})
			})
		})

		Convey(`Inspect`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("inspect")
				So(p, should.NotBeNil)
				p2 := app.App.Command("i")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call get and respond with template`, func() {
				app.App.Run([]string{"[app]", "inspect", "arrows.js"})
				So(a.GetCalledWith, should.Equal, "arrows.js")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"inspect", &api.AliasList{}})
			})
		})

		Convey(`Delete`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("delete")
				So(p, should.NotBeNil)
				p2 := app.App.Command("rm")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should prompt to delete and then confirm deleted`, func() {
				i.ReturnItem = true
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(a.DeleteCalledWith, should.Equal, "arrows.js")
				So(i.CallHistory[0], should.Resemble, []interface{}{"delete", "arrows.js"})
				So(i.CallHistory[1], should.Resemble, []interface{}{"deleted", "arrows.js"})
				i.ReturnItem = false
			})
			Convey(`Should prompt to delete and then confirm not deleted`, func() {
				i.ReturnItem = false
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(i.CallHistory[0], should.Resemble, []interface{}{"delete", "arrows.js"})
				So(i.CallHistory[1], should.Resemble, []interface{}{"notdeleted", "arrows.js"})
				i.ReturnItem = false
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
			Convey(`Should call get and respond with template`, func() {
				app.App.Run([]string{"[app]", "cat", "arrows.js"})
				So(a.GetCalledWith, should.Equal, "arrows.js")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"cat", &api.AliasList{}})
			})
		})

		Convey(`Rename`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("rename")
				So(p, should.NotBeNil)
				p2 := app.App.Command("mv")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call rename and respond with template`, func() {
				app.App.Run([]string{"[app]", "rename", "arrows.js", "pointers.js"})
				So(a.RenameCalledWith, should.Resemble, []string{"arrows.js", "pointers.js"})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"rename", &api.Alias{FullKey:"pointers.js"}})
			})
		})

		Convey(`Clone`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("clone")
				So(p, should.NotBeNil)
			})
			Convey(`Should call rename and respond with template`, func() {
				app.App.Run([]string{"[app]", "clone", "unicode/arrows.js", "myarrows.js"})
				So(a.CloneCalledWith, should.Resemble, []string{"unicode/arrows.js", "myarrows.js"})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"clone", &api.Alias{}})
			})
		})

		Convey(`Edit`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("edit")
				So(p, should.NotBeNil)
				p2 := app.App.Command("e")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call edit and respond with template`, func() {
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(o.EditCalledWith, should.Resemble, "arrows.js")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"edit", nil})
			})
			Convey(`Should call edit and respond with error if not exists`, func() {
				o.EditError = errors.New("Not found.")
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(o.EditCalledWith, should.Resemble, "arrows.js")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"edit", o.EditError})
				o.EditError = nil
			})
		})

		Convey(`Patch`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("patch")
				So(p, should.NotBeNil)
			})
			//can only patch with a fullKey, an ambiguous key will give 404
			//TODO: When updating a pinned kwklink, must force to give a new name
			// (since it is technically no longer the original)
			Convey(`Should call patch and respond with patch`, func() {
				app.App.Run([]string{"[app]", "patch", "arrows.js", "console.log('patched')"})
				So(a.PatchCalledWith, should.Resemble, []string{"arrows.js", "console.log('patched')"})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"patch", &api.Alias{FullKey:"arrows.js", Uri:"console.log('patched')"}})
			})
		})
	})
}


