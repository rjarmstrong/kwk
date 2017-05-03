package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/snippets"
	. "github.com/smartystreets/goconvey/convey"
	"bitbucket.com/sharingmachine/kwkcli/src/cmd"
	"testing"
)

func Test_Snippet(t *testing.T) {
	Convey("SNIPPET CLI", t, func() {
		app := CreateAppStub()
		ss := app.Snippets.(*snippets.ServiceMock)
		r := app.Runner.(*cmd.RunnerMock)
		//d := app.Dialogue.(*dlg.DialogMock)
		t := app.Settings.(*config.PersisterMock)
		//w := app.TemplateWriter.(*tmpl.WriterMock)
		//h := app.Search.(*search.TermMock)
		username := "richard"

		Convey(`Snippet Running`, func() {
			Convey(`Should call 'run' and open if found`, func() {
				a := models.NewAlias(username, "", "hola", "js")
				ss.ReturnItemsForGet = []*models.Snippet{{Alias: *a}}
				t.GetHydrates = &models.User{Username: a.Username}
				app.App.Run([]string{"[app]", "hola.js", "covert", "arg2"})
				So(ss.GetCalledWith, ShouldResemble, a)
				So(r.RunCalledWith, ShouldResemble, []interface{}{&ss.ReturnItemsForGet[0], []string{"covert", "arg2"}})
				ss.ReturnItemsForGet = nil
			})
			//	Convey(`Should call 'run' and prompt if multiple found`, func() {
			//		fullKey1 := "hola.sh"
			//		fullKey2 := "hola.js"
			//		a.ReturnItemsForGet = []*models.Snippet{
			//			{FullName: fullKey1},
			//			{FullName: fullKey2},
			//		}
			//		t.GetHydrates = &models.User{Username: "rjarmstrong"}
			//		d.MultiChoiceResponse = &dlg.DialogResponse{Value: a.ReturnItemsForGet[0]}
			//		app.App.Run([]string{"[app]", "hola", "arg1", "arg2"})
			//		So(a.GetCalledWith, ShouldResemble, &types.SnipName{Username: "rjarmstrong", FullKey: "hola"})
			//		So(d.MultiChoiceCalledWith, ShouldResemble, []interface{}{"dialog:choose", "Multiple matches. Choose a snippet to run:", a.ReturnItemsForGet})
			//		a.ReturnItemsForGet = nil
			//	})
			//	Convey(`Should suggest if not found`, func() {
			//		fullKey := "hola.sh"
			//		a.ReturnItemsForGet = []*models.Snippet{}
			//		results := []*models.SearchResult{}
			//		result := &models.SearchResult{Name:"suggestion"}
			//		results = append(results, result)
			//		h.ReturnForExecute = &models.SearchTermResponse{
			//			Total:1,
			//			Results:results,
			//		}
			//		app.App.Run([]string{"[app]", fullKey, "arg1", "arg2"})
			//		So(a.GetCalledWith, ShouldResemble, &types.SnipName{FullKey: "hola.sh"})
			//		So(w.RenderCalledWith[0], ShouldResemble, "search:alphaSuggest")
			//		So(w.RenderCalledWith[1].(*models.SearchTermResponse).Results[0], ShouldResemble, result)
			//		a.ReturnItemsForGet = nil
			//	})
			//})
			//
			//Convey(`New/Create/Save`, func() {
			//	Convey(`Create and Save should be equivalent`, func() {
			//		p := app.App.Command("new")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("create")
			//		p3 := app.App.Command("save")
			//		So(p2.Name, ShouldEqual, p.Name)
			//		So(p3.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should call create, save to clip board and respond with template WITH a fullKey`, func() {
			//		fullKey := "hola.sh"
			//		app.App.Run([]string{"[app]", "new", "echo hola", fullKey})
			//		So(a.CreateCalledWith, ShouldResemble, []string{"echo hola", fullKey})
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:new", &models.Snippet{FullName: fullKey}})
			//	})
			//	Convey(`Should call create, save to clip board and respond with template WITHOUT a fullKey`, func() {
			//		app.App.Run([]string{"[app]", "new", "echo hola"})
			//		So(a.CreateCalledWith, ShouldResemble, []string{"echo hola", ""})
			//		mockKey := "x5hi23"
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:new", &models.Snippet{FullName: mockKey}})
			//	})
			//})
			//
			//Convey(`Inspect`, func() {
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("inspect")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("i")
			//		So(p2.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should call get and respond with template`, func() {
			//		app.App.Run([]string{"[app]", "inspect", "arrows.js"})
			//		So(a.GetCalledWith, ShouldResemble, &types.SnipName{FullKey: "arrows.js"})
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:inspect", &models.SnippetList{}})
			//	})
			//})
			//
			//Convey(`Tag`, func() {
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("tag")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("t")
			//		So(p2.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should run by name (untag)`, func() {
			//		p := app.App.Command("untag")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("ut")
			//		So(p2.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should call tag and respond with template`, func() {
			//		app.App.Run([]string{"[app]", "tag", "arrows.js", "tag1", "tag2"})
			//		So(a.TagCalledWith, ShouldResemble, map[string][]string{
			//			"arrows.js": []string{"tag1", "tag2"},
			//		})
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:tag", &models.Snippet{}})
			//	})
			//	Convey(`Should call untag and respond with template`, func() {
			//		app.App.Run([]string{"[app]", "untag", "arrows.js", "tag1", "tag2"})
			//		So(a.UnTagCalledWith, ShouldResemble, map[string][]string{
			//			"arrows.js": []string{"tag1", "tag2"},
			//		})
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:untag", &models.Snippet{}})
			//	})
			//})
			//
			//Convey(`Delete`, func() {
			//	data := &models.Snippet{FullName: "arrows.js"}
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("delete")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("rm")
			//		So(p2.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should prompt to delete and then confirm deleted`, func() {
			//		d.ReturnItem = &dlg.DialogResponse{Ok: true}
			//		app.App.Run([]string{"[app]", "delete", "arrows.js"})
			//		So(a.DeleteCalledWith, ShouldEqual, "arrows.js")
			//		So(d.LastModalCalledWith[0].(string), ShouldResemble, "snippet:delete")
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:deleted", data})
			//		d.ReturnItem = nil
			//	})
			//	Convey(`Should prompt to delete and then confirm not deleted`, func() {
			//		d.ReturnItem = &dlg.DialogResponse{Ok: false}
			//		app.App.Run([]string{"[app]", "delete", "arrows.js"})
			//		So(d.CallHistory[0].([]interface{})[0], ShouldResemble, "snippet:delete")
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:notdeleted", data})
			//		d.ReturnItem = nil
			//	})
			//})
			//
			//Convey(`Cat`, func() {
			//	Convey(`Should test equivalent cmd names`, func() {
			//		p := app.App.Command("cat")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("raw")
			//		p3 := app.App.Command("read")
			//		p4 := app.App.Command("print")
			//		p5 := app.App.Command("get")
			//		So(p2.Name, ShouldEqual, p.Name)
			//		So(p3.Name, ShouldEqual, p.Name)
			//		So(p4.Name, ShouldEqual, p.Name)
			//		So(p5.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should call get and respond with template`, func() {
			//		a.ReturnItemsForGet = []*models.Snippet{
			//			{FullName: "arrows.js"},
			//		}
			//		app.App.Run([]string{"[app]", "cat", "arrows.js"})
			//		So(a.GetCalledWith, ShouldResemble, &types.SnipName{FullKey: "arrows.js"})
			//		So(w.RenderCalledWith[0].(string), ShouldResemble, "snippet:cat")
			//		So(w.RenderCalledWith[1].(models.Snippet).FullName, ShouldResemble, "arrows.js")
			//	})
			//})
			//
			//Convey(`Rename`, func() {
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("rename")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("mv")
			//		So(p2.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should call rename and respond with template`, func() {
			//		app.App.Run([]string{"[app]", "rename", "arrows.js", "pointers.js"})
			//		So(a.RenameCalledWith, ShouldResemble, []string{"arrows.js", "pointers.js"})
			//		So(w.RenderCalledWith[0].(string), ShouldResemble, "snippet:renamed")
			//		r := w.RenderCalledWith[1].(*map[string]string)
			//		So((*r)["newFullKey"], ShouldResemble,  "pointers.js")
			//	})
			//})
			//
			//Convey(`Clone`, func() {
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("clone")
			//		So(p, ShouldNotBeNil)
			//	})
			//	Convey(`Should call rename and respond with template`, func() {
			//		app.App.Run([]string{"[app]", "clone", "unicode/arrows.js", "myarrows.js"})
			//		So(a.CloneCalledWith, ShouldResemble, []interface{}{
			//
			//			&types.SnipName{Username: "unicode", FullKey: "arrows.js"}, "myarrows.js"})
			//
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:cloned", &models.Snippet{}})
			//	})
			//})
			//
			//Convey(`Edit`, func() {
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("edit")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("e")
			//		So(p2.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should call edit and respond with template`, func() {
			//		a.ReturnItemsForGet = []*models.Snippet{{FullName: "arrows.js"}}
			//		d.MultiChoiceResponse = &dlg.DialogResponse{Value: a.ReturnItemsForGet[0]}
			//		app.App.Run([]string{"[app]", "edit", "arrows.js"})
			//		So(r.EditCalledWith, ShouldResemble, &a.ReturnItemsForGet[0])
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:edited", &a.ReturnItemsForGet[0]})
			//	})
			//	Convey(`Should call edit and respond with error if not exists`, func() {
			//		app.App.Run([]string{"[app]", "edit", "arrows.js"})
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:notfound", &models.Snippet{FullName: "arrows.js"}})
			//	})
			//})
			//
			//Convey(`Patch`, func() {
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("patch")
			//		So(p, ShouldNotBeNil)
			//	})
			//	Convey(`Should call patch and respond with patch`, func() {
			//		app.App.Run([]string{"[app]", "patch", "arrows.js", "console.log('original')", "console.log('patched')"})
			//		So(a.PatchCalledWith, ShouldResemble, []string{"arrows.js", "console.log('original')", "console.log('patched')"})
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:patched", &models.Snippet{FullName: "arrows.js", Snip: "console.log('patched')"}})
			//	})
			//})
			//
			//Convey(`List`, func() {
			//	Convey(`Should run by name`, func() {
			//		p := app.App.Command("list")
			//		So(p, ShouldNotBeNil)
			//		p2 := app.App.Command("ls")
			//		So(p2.Name, ShouldEqual, p.Name)
			//	})
			//	Convey(`Should call list and respond with template`, func() {
			//		app.App.Run([]string{"[app]", "list", "5", "tag1"})
			//		So(a.ListCalledWith.Username, ShouldResemble, "")
			//		So(a.ListCalledWith.Size, ShouldResemble, int64(5))
			//		So(a.ListCalledWith.Since/1000, ShouldResemble, time.Now().UnixNano()/1000000000)
			//		So(a.ListCalledWith.Tags, ShouldResemble, []string{"tag1"})
			//		So(w.RenderCalledWith, ShouldResemble, []interface{}{"snippet:list", &models.SnippetList{}})
			//	})
			//})
		})
	})
}