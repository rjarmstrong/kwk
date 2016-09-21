package routes

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/openers"
	_ "github.com/iris-contrib/errors"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
	"github.com/kwk-links/kwk-cli/libs/models"
	"github.com/kwk-links/kwk-cli/libs/app"
)

func Test_Alias(t *testing.T) {
	Convey("ALIAS COMMANDS", t, func() {
		a := &aliases.AliasesMock{}
		i := &gui.InteractionMock{}
		s := &system.SystemMock{}
		o := &openers.OpenerMock{}
		sett := &settings.SettingsMock{}
		appl := app.NewKwkApp(a, s, sett, i, o)

		Convey(`Command not found`, func() {
			Convey(`Should call get and open if found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Alias{
					{FullKey:fullKey},
				}
				appl.App.Run([]string{"[app]", fullKey, "covert", "arg2"})
				So(a.GetCalledWith, should.Equal, fullKey)
				So(o.OpenCalledWith, should.Resemble, []interface{}{&a.ReturnItemsForGet[0], []string{"covert", "arg2"}})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should call 'get' and prompt if multiple found`, func() {
				fullKey1 := "hola.sh"
				fullKey2 := "hola.js"
				a.ReturnItemsForGet = []models.Alias{
					{FullKey:fullKey1},
					{FullKey:fullKey2},
				}
				appl.App.Run([]string{"[app]", "hola", "arg1", "arg2"})
				So(a.GetCalledWith, should.Equal, "hola")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"chooseBetweenKeys", a.ReturnItemsForGet})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should respond if not found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Alias{}
				appl.App.Run([]string{"[app]", fullKey, "arg1", "arg2"})
				So(a.GetCalledWith, should.Equal, fullKey)
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"notfound", fullKey})
				a.ReturnItemsForGet = nil
			})
		})

		Convey(`New`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("new")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("create")
				p3 := appl.App.Command("save")
				So(p2.Name, should.Equal, p.Name)
				So(p3.Name, should.Equal, p.Name)
			})
			Convey(`Should call create, save to clip board and respond with template WITH a fullKey`, func() {
				fullKey := "hola.sh"
				appl.App.Run([]string{"[app]", "new", "echo hola", fullKey})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", fullKey})
				So(s.CopyToClipboardCalledWith, should.Equal, fullKey)
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"new", &models.Alias{FullKey:fullKey}})
			})
			Convey(`Should call create, save to clip board and respond with template WITHOUT a fullKey`, func() {
				appl.App.Run([]string{"[app]", "new", "echo hola"})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", ""})
				mockKey := "x5hi23"
				So(s.CopyToClipboardCalledWith, should.Equal, mockKey)
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"new", &models.Alias{FullKey:mockKey}})
			})
		})

		Convey(`Inspect`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("inspect")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("i")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call get and respond with template`, func() {
				appl.App.Run([]string{"[app]", "inspect", "arrows.js"})
				So(a.GetCalledWith, should.Equal, "arrows.js")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"inspect", &models.AliasList{}})
			})
		})

		Convey(`Tag`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("tag")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("t")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should run by name (untag)`, func() {
				p := appl.App.Command("untag")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("ut")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call tag and respond with template`, func() {
				appl.App.Run([]string{"[app]", "tag", "arrows.js", "tag1", "tag2"})
				So(a.TagCalledWith, should.Resemble, map[string][]string {
					"arrows.js" : []string{"tag1", "tag2"},
				})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"tag", &models.Alias{}})
			})
			Convey(`Should call untag and respond with template`, func() {
				appl.App.Run([]string{"[app]", "untag", "arrows.js", "tag1", "tag2"})
				So(a.UnTagCalledWith, should.Resemble, map[string][]string {
					"arrows.js" : []string{"tag1", "tag2"},
				})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"untag", &models.Alias{}})
			})
		})

		Convey(`Delete`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("delete")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("rm")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should prompt to delete and then confirm deleted`, func() {
				i.ReturnItem = true
				appl.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(a.DeleteCalledWith, should.Equal, "arrows.js")
				So(i.CallHistory[0], should.Resemble, []interface{}{"delete", "arrows.js"})
				So(i.CallHistory[1], should.Resemble, []interface{}{"deleted", "arrows.js"})
				i.ReturnItem = false
			})
			Convey(`Should prompt to delete and then confirm not deleted`, func() {
				i.ReturnItem = false
				appl.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(i.CallHistory[0], should.Resemble, []interface{}{"delete", "arrows.js"})
				So(i.CallHistory[1], should.Resemble, []interface{}{"notdeleted", "arrows.js"})
				i.ReturnItem = false
			})
		})

		Convey(`Cat`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("cat")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("raw")
				p3 := appl.App.Command("read")
				p4 := appl.App.Command("print")
				p5 := appl.App.Command("get")
				So(p2.Name, should.Equal, p.Name)
				So(p3.Name, should.Equal, p.Name)
				So(p4.Name, should.Equal, p.Name)
				So(p5.Name, should.Equal, p.Name)
			})
			Convey(`Should call get and respond with template`, func() {
				appl.App.Run([]string{"[app]", "cat", "arrows.js"})
				So(a.GetCalledWith, should.Equal, "arrows.js")
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"cat", &models.AliasList{}})
			})
		})

		Convey(`Rename`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("rename")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("mv")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call rename and respond with template`, func() {
				appl.App.Run([]string{"[app]", "rename", "arrows.js", "pointers.js"})
				So(a.RenameCalledWith, should.Resemble, []string{"arrows.js", "pointers.js"})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"rename", &models.Alias{FullKey:"pointers.js"}})
			})
		})

		Convey(`Clone`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("clone")
				So(p, should.NotBeNil)
			})
			Convey(`Should call rename and respond with template`, func() {
				appl.App.Run([]string{"[app]", "clone", "unicode/arrows.js", "myarrows.js"})
				So(a.CloneCalledWith, should.Resemble, []string{"unicode/arrows.js", "myarrows.js"})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"clone", &models.Alias{}})
			})
		})

		Convey(`Edit`, func() {
			//Convey(`Should run by name`, func() {
			//	p := app.App.Command("edit")
			//	So(p, should.NotBeNil)
			//	p2 := app.App.Command("e")
			//	So(p2.Name, should.Equal, p.Name)
			//})
			//Convey(`Should call edit and respond with template`, func() {
			//	a.ReturnItemsForGet = []models.Alias{
			//			{FullKey:"arrows.js"}}
			//	app.App.Run([]string{"[app]", "edit", "arrows.js"})
			//	So(o.EditCalledWith, should.Resemble, &a.ReturnItemsForGet[0])
			//	So(i.LastRespondCalledWith, should.Resemble, []interface{}{"edit", nil})
			//})
			//Convey(`Should call edit and respond with error if not exists`, func() {
			//	o.EditError = errors.New("Not found.")
			//	app.App.Run([]string{"[app]", "edit", "arrows.js"})
			//	//So(o.EditCalledWith, should.Resemble, "arrows.js")
			//	So(i.LastRespondCalledWith, should.Resemble, []interface{}{"edit", o.EditError})
			//	o.EditError = nil
			//})
		})

		Convey(`Patch`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("patch")
				So(p, should.NotBeNil)
			})
			//can only patch with a fullKey, an ambiguous key will give 404
			//TODO: When updating a pinned kwklink, must force to give a new name
			// (since it is technically no longer the original)
			Convey(`Should call patch and respond with patch`, func() {
				appl.App.Run([]string{"[app]", "patch", "arrows.js", "console.log('patched')"})
				So(a.PatchCalledWith, should.Resemble, []string{"arrows.js", "console.log('patched')"})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"patch", &models.Alias{FullKey:"arrows.js", Uri:"console.log('patched')"}})
			})
		})

		Convey(`List`, func() {
			Convey(`Should run by name`, func() {
				p := appl.App.Command("list")
				So(p, should.NotBeNil)
				p2 := appl.App.Command("ls")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call list and respond with template`, func() {
				appl.App.Run([]string{"[app]", "list", "3", "5"})
				So(a.ListCalledWith, should.Resemble, []string{"3", "5"})
				So(i.LastRespondCalledWith, should.Resemble, []interface{}{"list", &models.AliasList{}})
			})
		})
	})
}


