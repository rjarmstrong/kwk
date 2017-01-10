package app

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"testing"
	"time"
)

func Test_Snippet(t *testing.T) {
	Convey("SNIPPET CLI", t, func() {
		app := CreateAppStub()
		a := app.Snippets.(*snippets.ServiceMock)
		r := app.Runner.(*cmd.RunnerMock)
		s := app.System.(*sys.ManagerMock)
		d := app.Dialogue.(*dlg.DialogMock)
		t := app.Settings.(*config.SettingsMock)
		w := app.TemplateWriter.(*tmpl.WriterMock)
		h := app.Search.(*search.TermMock)

		Convey(`Snippet Running`, func() {
			Convey(`Should call 'run' and open if found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Snippet{
					{FullName: fullKey},
				}
				t.GetHydrateWith = &models.User{Username: "rjarmstrong"}
				app.App.Run([]string{"[app]", fullKey, "covert", "arg2"})
				So(a.GetCalledWith, should.Resemble, &models.Alias{FullKey: fullKey, Username: "rjarmstrong"})
				So(r.OpenCalledWith, should.Resemble, []interface{}{&a.ReturnItemsForGet[0], []string{"covert", "arg2"}})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should call 'run' and prompt if multiple found`, func() {
				fullKey1 := "hola.sh"
				fullKey2 := "hola.js"
				a.ReturnItemsForGet = []models.Snippet{
					{FullName: fullKey1},
					{FullName: fullKey2},
				}
				t.GetHydrateWith = &models.User{Username: "rjarmstrong"}
				d.MultiChoiceResponse = &dlg.DialogResponse{Value: a.ReturnItemsForGet[0]}
				app.App.Run([]string{"[app]", "hola", "arg1", "arg2"})
				So(a.GetCalledWith, should.Resemble, &models.Alias{Username: "rjarmstrong", FullKey: "hola"})
				So(d.MultiChoiceCalledWith, should.Resemble, []interface{}{"dialog:choose", "Multiple matches. Choose a snippet to run:", a.ReturnItemsForGet})
				a.ReturnItemsForGet = nil
			})
			Convey(`Should suggest if not found`, func() {
				fullKey := "hola.sh"
				a.ReturnItemsForGet = []models.Snippet{}
				results := []*models.SearchResult{}
				result := &models.SearchResult{Key:"suggestion"}
				results = append(results, result)
				h.ReturnForExecute = &models.SearchTermResponse{
					Total:1,
					Results:results,
				}
				app.App.Run([]string{"[app]", fullKey, "arg1", "arg2"})
				So(a.GetCalledWith, should.Resemble, &models.Alias{FullKey: "hola.sh"})
				So(w.RenderCalledWith[0], should.Resemble, "search:alphaSuggest")
				So(w.RenderCalledWith[1].(*models.SearchTermResponse).Results[0], should.Resemble, result)
				a.ReturnItemsForGet = nil
			})
		})

		Convey(`New/Create/Save`, func() {
			Convey(`Create and Save should be equivalent`, func() {
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:new", &models.Snippet{FullName: fullKey}})
			})
			Convey(`Should call create, save to clip board and respond with template WITHOUT a fullKey`, func() {
				app.App.Run([]string{"[app]", "new", "echo hola"})
				So(a.CreateCalledWith, should.Resemble, []string{"echo hola", ""})
				mockKey := "x5hi23"
				So(s.CopyToClipboardCalledWith, should.Equal, mockKey)
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:new", &models.Snippet{FullName: mockKey}})
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
				So(a.GetCalledWith, should.Resemble, &models.Alias{FullKey: "arrows.js"})
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
			data := &models.Snippet{FullName: "arrows.js"}
			Convey(`Should run by name`, func() {
				p := app.App.Command("delete")
				So(p, should.NotBeNil)
				p2 := app.App.Command("rm")
				So(p2.Name, should.Equal, p.Name)
			})
			Convey(`Should prompt to delete and then confirm deleted`, func() {
				d.ReturnItem = &dlg.DialogResponse{Ok: true}
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(a.DeleteCalledWith, should.Equal, "arrows.js")
				So(d.LastModalCalledWith[0].(string), should.Resemble, "snippet:delete")
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:deleted", data})
				d.ReturnItem = nil
			})
			Convey(`Should prompt to delete and then confirm not deleted`, func() {
				d.ReturnItem = &dlg.DialogResponse{Ok: false}
				app.App.Run([]string{"[app]", "delete", "arrows.js"})
				So(d.CallHistory[0].([]interface{})[0], should.Resemble, "snippet:delete")
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:notdeleted", data})
				d.ReturnItem = nil
			})
		})

		Convey(`Cat`, func() {
			Convey(`Should test equivalent cmd names`, func() {
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
					{FullName: "arrows.js"},
				}
				app.App.Run([]string{"[app]", "cat", "arrows.js"})
				So(a.GetCalledWith, should.Resemble, &models.Alias{FullKey: "arrows.js"})
				So(w.RenderCalledWith[0].(string), should.Resemble, "snippet:cat")
				So(w.RenderCalledWith[1].(models.Snippet).FullName, should.Resemble, "arrows.js")
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

					&models.Alias{Username: "unicode", FullKey: "arrows.js"}, "myarrows.js"})

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
				a.ReturnItemsForGet = []models.Snippet{{FullName: "arrows.js"}}
				d.MultiChoiceResponse = &dlg.DialogResponse{Value: a.ReturnItemsForGet[0]}
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(r.EditCalledWith, should.Resemble, &a.ReturnItemsForGet[0])
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:edited", &a.ReturnItemsForGet[0]})
			})
			Convey(`Should call edit and respond with error if not exists`, func() {
				app.App.Run([]string{"[app]", "edit", "arrows.js"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:notfound", &models.Snippet{FullName: "arrows.js"}})
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
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:patched", &models.Snippet{FullName: "arrows.js", Snip: "console.log('patched')"}})
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
				So(a.ListCalledWith[0], should.Resemble, "")
				So(a.ListCalledWith[1], should.Resemble, int64(5))
				So(a.ListCalledWith[2].(int64)/1000, should.Resemble, time.Now().UnixNano()/1000000000)
				So(a.ListCalledWith[3], should.Resemble, []string{"tag1"})
				So(w.RenderCalledWith, should.Resemble, []interface{}{"snippet:list", &models.SnippetList{}})
			})
		})
	})
}
