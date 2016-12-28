package app

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/openers"
	"bitbucket.com/sharingmachine/kwkcli/settings"
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_Snippet(t *testing.T) {
	Convey("SNIPPET CLI", t, func() {
		app := CreateAppStub()
		a := app.Snippets.(*snippets.Mock)
		o := app.Openers.(*openers.OpenerMock)
		s := app.System.(*system.MockSystem)
		d := app.Dialogues.(*dlg.MockDialogue)
		t := app.Settings.(*settings.SettingsMock)
		w := app.TemplateWriter.(*tmpl.MockWriter)

		Convey(`Command not found`, func() {
			Convey(`Should call get and open if found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Snippet{
					{FullKey: fullKey},
				}
				t.GetHydrateWith = &models.User{Username: "rjarmstrong"}
				app.App.Run([]string{"[app]", fullKey, "covert", "arg2"})
				So(a.GetCalledWith, should.Resemble, &models.KwkKey{FullKey: fullKey, Username: "rjarmstrong"})
				So(o.OpenCalledWith, should.Resemble, []interface{}{&a.ReturnItemsForGet[0], []string{"covert", "arg2"}})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should call 'get' and prompt if multiple found`, func() {
				fullKey1 := "hola.sh"
				fullKey2 := "hola.js"
				a.ReturnItemsForGet = []models.Snippet{
					{FullKey: fullKey1},
					{FullKey: fullKey2},
				}
				t.GetHydrateWith = &models.User{Username: "rjarmstrong"}
				d.MultiChoiceResponse = &dlg.DialogueResponse{Value: &a.ReturnItemsForGet[0]}
				app.App.Run([]string{"[app]", "hola", "arg1", "arg2"})
				So(a.GetCalledWith, should.Resemble, &models.KwkKey{Username: "rjarmstrong", FullKey: "hola"})
				So(d.MultiChoiceCalledWith, should.Resemble, []interface{}{"snippet:choose", nil,
					a.ReturnItemsForGet})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should respond if not found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Snippet{}
				app.App.Run([]string{"[app]", fullKey, "arg1", "arg2"})
				So(a.GetCalledWith, should.Resemble, &models.KwkKey{FullKey: "hola.sh"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:notfound",
					map[string]interface{}{
						"fullKey": "hola.sh",
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:new", &models.Snippet{FullKey: fullKey}})
			})
			Convey(`Should call create, save to clip board and respond with template WITHOUT a fullKey`, func() {
				app.App.Run([]string{"[app]", "new", "echo hola"})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", ""})
				mockKey := "x5hi23"
				So(s.CopyToClipboardCalledWith, should.Equal, mockKey)
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:new", &models.Snippet{FullKey: mockKey}})
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
				So(a.GetCalledWith, should.Resemble, &models.KwkKey{FullKey: "arrows.js"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:inspect", &models.SnippetList{}})
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
				So(a.TagCalledWith, should.Resemble, map[string][]string{
					"arrows.js": []string{"tag1", "tag2"},
				})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:tag", &models.Snippet{}})
			})
			Convey(`Should call untag and respond with template`, func() {
				app.App.Run([]string{"[app]", "untag", "arrows.js", "tag1", "tag2"})
				So(a.UnTagCalledWith, should.Resemble, map[string][]string{
					"arrows.js": []string{"tag1", "tag2"},
				})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:untag", &models.Snippet{}})
			})
		})

		Convey(`Delete`, func() {
			data := &models.Snippet{FullKey: "arrows.js"}
			Convey(`Should run by name`, func() {
				p := app.App.Command("delete")
				So(p, should.NotBeNil)
				p2 := app.App.Command("rm")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should prompt to delete and then confirm deleted`, func() {
				d.ReturnItem = &dlg.DialogueResponse{Ok: true}
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(a.DeleteCalledWith, should.Equal, "arrows.js")
				So(d.LastModalCalledWith[0].(string), should.Resemble, "snippet:delete")
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:deleted", data})
				d.ReturnItem = nil
			})
			Convey(`Should prompt to delete and then confirm not deleted`, func() {
				d.ReturnItem = &dlg.DialogueResponse{Ok: false}
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(d.CallHistory[0].([]interface{})[0], should.Resemble, "snippet:delete")
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:notdeleted", data})
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
				a.ReturnItemsForGet = []models.Snippet{
					{FullKey: "arrows.js"},
				}
				app.App.Run([]string{"[app]", "cat", "arrows.js"})
				So(a.GetCalledWith, should.Resemble, &models.KwkKey{FullKey: "arrows.js"})
				So(w.RenderCalledWith[0].(string), should.Resemble, "snippet:cat")
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
				So(w.RenderCalledWith[0].(string), should.Resemble, "snippet:renamed")
				r := w.RenderCalledWith[1].(*map[string]string)
				So((*r)["newFullKey"], should.Resemble,  "pointers.js")
			})
		})

		Convey(`Clone`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("clone")
				So(p, should.NotBeNil)
			})
			Convey(`Should call rename and respond with template`, func() {
				app.App.Run([]string{"[app]", "clone", "unicode/arrows.js", "myarrows.js"})
				So(a.CloneCalledWith, should.Resemble, []interface{}{

					&models.KwkKey{Username: "unicode", FullKey: "arrows.js"}, "myarrows.js"})

				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:cloned", &models.Snippet{}})
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
				a.ReturnItemsForGet = []models.Snippet{{FullKey: "arrows.js"}}
				d.MultiChoiceResponse = &dlg.DialogueResponse{Value: a.ReturnItemsForGet[0]}
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(o.EditCalledWith, should.Resemble, &a.ReturnItemsForGet[0])
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:edited", &a.ReturnItemsForGet[0]})
			})
			Convey(`Should call edit and respond with error if not exists`, func() {
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:notfound", &models.Snippet{FullKey: "arrows.js"}})
			})
		})

		Convey(`Patch`, func() {
			Convey(`Should run by name`, func() {
				p := app.App.Command("patch")
				So(p, should.NotBeNil)
			})
			Convey(`Should call patch and respond with patch`, func() {
				app.App.Run([]string{"[app]", "patch", "arrows.js", "console.log('original')", "console.log('patched')"})
				So(a.PatchCalledWith, should.Resemble, []string{"arrows.js", "console.log('original')", "console.log('patched')"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:patched", &models.Snippet{FullKey: "arrows.js", Snip: "console.log('patched')"}})
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
				app.App.Run([]string{"[app]", "list", "5", "tag1"})
				So(a.ListCalledWith, should.Resemble, []interface{}{"", int64(5), time.Now().UnixNano()/1000000*1000, []string{"tag1"}})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:list", &models.SnippetList{}})
			})
		})
	})
}
