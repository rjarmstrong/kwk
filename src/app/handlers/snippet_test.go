package handlers

import (
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

var snippet1 = response{val: &types.ListResponse{
	Suggested: false,
	Items: []*types.Snippet{
		{
			Alias: types.NewAlias("", "pouch1", "snippet1", "txt"),
		},
	}}}

var snippet1Suggested = response{val: &types.ListResponse{
	Suggested: true,
	Items: []*types.Snippet{
		{
			Alias: types.NewAlias("", "pouch1", "snippet1", "txt"),
		},
	}}}

var johnnyRoot = response{val: &types.RootResponse{
	Username: "johnny",
	Pouches: []*types.Pouch{
		{Name: "pouch1"},
	},
}}

func TestSnippets_Create(t *testing.T) {
	var cases = []struct {
		args    []string
		pipe    bool
		content string
		code    errs.ErrCode
	}{
		{args: []string{"content"}, pipe: false, code: errs.CodeInvalidArgument},
		{args: []string{"conte$nt"}, pipe: false, code: errs.CodeInvalidArgument},
		{args: []string{"pouch/name"}, pipe: false, content: ""},
		{args: []string{"pouch/name", "content"}, pipe: false, content: "content"},
		{args: []string{"content", "pouch/name"}, pipe: false, content: "content"},
		{args: []string{"content", "name"}, pipe: false, code: errs.CodeInvalidArgument},
		{args: []string{"pouch/name", "pouch/name"}, pipe: false, code: errs.CodeInvalidArgument},
		{args: []string{"pouch$)/name", "content"}, pipe: false, code: errs.CodeInvalidArgument},
	}

	for _, c := range cases {
		err := snippets.Create(c.args, c.pipe)
		if err != nil {
			assert.True(t, errs.HasCode(err, c.code), "Case: %+v %+v", err.(*errs.Error).Code, c)
			snippetClient.PopCalled("Create")
			continue
		}
		requestContent := snippetClient.PopCalled("Create").(*types.CreateRequest).Content
		assert.Equal(t, c.content, requestContent, "Case: %+v", c)
	}
}

func TestSnippets_Search(t *testing.T) {
	prefs.PrivateView = true
	un := "username1"
	err := snippets.Search(un, "term")
	assert.Nil(t, err)
	req := snippetClient.PopCalled("Alpha").(*types.AlphaRequest)
	assert.Equal(t, true, req.PrivateView)
	assert.Equal(t, un, req.Username)

	prefs.PrivateView = false
	err = snippets.Search(un, "term")
	assert.Nil(t, err)
	req = snippetClient.PopCalled("Alpha").(*types.AlphaRequest)
	assert.Equal(t, false, req.PrivateView)
}

func TestSnippets_ViewListOrRun(t *testing.T) {

	t.Log("VIEW a snippet")
	snippetClient.returns["GetRoot"] = johnnyRoot
	snippetClient.returns["Get"] = snippet1
	err := snippets.ViewListOrRun("name1", true)
	assert.Nil(t, err)
	handler := writer.PopCalled("EWrite")
	assert.Equal(t, "SnippetView", handler)
	uri := snippetClient.PopCalled("Get").(*types.GetRequest).Alias.URI()
	assert.Equal(t, "name1", uri)

	t.Log("RUN a snippet")
	err = snippets.ViewListOrRun("name1", false, "y", "z")
	assert.Nil(t, err)
	uri = runner.PopCalled("Run").(string)
	assert.Equal(t, "pouch1/snippet1.txt", uri)

	t.Log("LIST a pouch")
	err = snippets.ViewListOrRun("pouch1", true)
	assert.Nil(t, err)
	lreq := snippetClient.PopCalled("List").(*types.ListRequest)
	assert.Equal(t, "pouch1", lreq.Pouch)

	t.Log("LIST root")
	err = snippets.ViewListOrRun("/johnny", true)
	assert.Nil(t, err)
	assert.Equal(t, "johnny", rootCalled.Username)
	rootCalled = nil
}

func TestSnippets_Cat(t *testing.T) {
	t.Log("EXACT MATCH")
	snippetClient.returns["Get"] = snippet1
	err := snippets.Cat("name1")
	assert.Nil(t, err)
	handler := writer.PopCalled("EWrite")
	assert.Equal(t, "SnippetCat", handler)

	t.Log("SUGGEST")
	snippetClient.returns["Get"] = snippet1Suggested
	err = snippets.Cat("name1")
	assert.Nil(t, err)
	funcName := dlg.PopCalled("Modal")
	assert.Equal(t, "Info", funcName)
	handler = writer.PopCalled("EWrite")
	assert.Equal(t, "SnippetCat", handler)
}

//func TestSnippets_Run(t *testing.T) {
//	snippets.Run()
//}
