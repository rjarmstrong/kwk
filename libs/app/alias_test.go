package app

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
)


func Test_Alias(t *testing.T) {
	Convey("ALIAS COMMANDS", t, func() {
		app := CreateAppStub()
		a := app.Aliases.(*aliases.AliasesMock)
		o := app.Openers.(*openers.OpenerMock)
		s := app.System.(*system.SystemMock)
		d := app.Dialogues.(*gui.DialogueMock)
		w := app.TemplateWriter.(*gui.TemplateWriterMock)

		Convey(`Command not found`, func() {
			Convey(`Should call get and open if found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Alias{
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
				a.ReturnItemsForGet = []models.Alias{
					{FullKey:fullKey1},
					{FullKey:fullKey2},
				}
				d.MultiChoiceResponse = &gui.DialogueResponse{Value:&a.ReturnItemsForGet[0]}
				app.App.Run([]string{"[app]", "hola", "arg1", "arg2"})
				So(a.GetCalledWith, should.Equal, "hola")
				So(d.MultiChoiceCalledWith, should.Resemble, []interface{}{"alias:choose", nil,
					[]interface{}{a.ReturnItemsForGet}})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should respond if not found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Alias{}
				app.App.Run([]string{"[app]", fullKey, "arg1", "arg2"})
				So(a.GetCalledWith, should.Equal, fullKey)
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:notfound",
					map[string]interface{}{
						"fullKey" : "hola.sh",
					}})
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:new", &models.Alias{FullKey:fullKey}})
			})
			Convey(`Should call create, save to clip board and respond with template WITHOUT a fullKey`, func() {
				app.App.Run([]string{"[app]", "new", "echo hola"})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", ""})
				mockKey := "x5hi23"
				So(s.CopyToClipboardCalledWith, should.Equal, mockKey)
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:new", &models.Alias{FullKey:mockKey}})
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:inspect", &models.AliasList{}})
			})
		})

		Convey(`Tag`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("tag")
				So(p, should.NotBeNil)
				p2 := app.App.Command("t")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should run by name (untag)`, func() {
				p := app.App.Command("untag")
				So(p, should.NotBeNil)
				p2 := app.App.Command("ut")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call tag and respond with template`, func() {
				app.App.Run([]string{"[app]", "tag", "arrows.js", "tag1", "tag2"})
				So(a.TagCalledWith, should.Resemble, map[string][]string {
					"arrows.js" : []string{"tag1", "tag2"},
				})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:tag", &models.Alias{}})
			})
			Convey(`Should call untag and respond with template`, func() {
				app.App.Run([]string{"[app]", "untag", "arrows.js", "tag1", "tag2"})
				So(a.UnTagCalledWith, should.Resemble, map[string][]string {
					"arrows.js" : []string{"tag1", "tag2"},
				})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:untag", &models.Alias{}})
			})
		})

		Convey(`Delete`, func() {
			data := map[string]string{"fullKey" : "arrows.js"}
			deletePrompt := []interface{}{"alias:delete", data}
			Convey(`Should run by name`, func() {
				p := app.App.Command("delete")
				So(p, should.NotBeNil)
				p2 := app.App.Command("rm")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should prompt to delete and then confirm deleted`, func() {
				d.ReturnItem = &gui.DialogueResponse{Ok:true}
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(a.DeleteCalledWith, should.Equal, "arrows.js")
				So(d.LastModalCalledWith, should.Resemble, deletePrompt)
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:deleted", map[string]string{"fullKey":"arrows.js"}})
				d.ReturnItem = nil
			})
			Convey(`Should prompt to delete and then confirm not deleted`, func() {
				d.ReturnItem = &gui.DialogueResponse{Ok:false}
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(d.CallHistory[0], should.Resemble, deletePrompt)
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:notdeleted", data})
				d.ReturnItem = nil
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:cat", &models.AliasList{}})
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:renamed", &models.Alias{FullKey:"pointers.js"}})
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:cloned", &models.Alias{}})
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
				a.ReturnItemsForGet = []models.Alias{{FullKey:"arrows.js"}}
				d.MultiChoiceResponse = &gui.DialogueResponse{Value:a.ReturnItemsForGet[0]}
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(o.EditCalledWith, should.Resemble, &a.ReturnItemsForGet[0])
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:edited", &a.ReturnItemsForGet[0]})
			})
			Convey(`Should call edit and respond with error if not exists`, func() {
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:notfound", map[string]interface{}{"fullKey":"arrows.js"}})
			})
		})

		Convey(`Patch`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("patch")
				So(p, should.NotBeNil)
			})
			//can only patch with a fullKey, an ambiguous key will give 404
			//TODO: When updating a pinned alias, must force to give a new name
			// (since it is technically no longer the original)
			Convey(`Should call patch and respond with patch`, func() {
				app.App.Run([]string{"[app]", "patch", "arrows.js", "console.log('patched')"})
				So(a.PatchCalledWith, should.Resemble, []string{"arrows.js", "console.log('patched')"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:patched", &models.Alias{FullKey:"arrows.js", Uri:"console.log('patched')"}})
			})
		})

		Convey(`List`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("list")
				So(p, should.NotBeNil)
				p2 := app.App.Command("ls")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should call list and respond with template`, func() {
				app.App.Run([]string{"[app]", "list", "3", "5", "tag1"})
				So(a.ListCalledWith, should.Resemble, []interface{}{"richard", int32(3), int32(5), []string{"tag1"}})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"alias:list", &models.AliasList{}})
			})
		})
	})
}


