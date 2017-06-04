package runtime

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/types"
	"github.com/stretchr/testify/assert"
)

func Test_Editor(t *testing.T) {
	out.DebugEnabled = true
	fileMock := &snippetReadWriter{}
	//username := "test-man"
	env := &yaml.MapSlice{}
	prefs := &out.Prefs{}
	yaml.Unmarshal([]byte(testEnvString), env)
	var patcherCalled int
	var content = "some content"
	fileMock.readVal = content
	editor := NewEditor(env, prefs, patcher(&patcherCalled), fileMock).(*editor)
	editor.guiFunc = mockRunner
	s := types.NewBlankSnippet()
	s.Content = content

	t.Log("If file unchanged should Not patch file.")
	err := editor.Edit(s)
	assert.Nil(t, err)
	assert.Equal(t, 0, patcherCalled)

	t.Log("If file changed Should patch file.")
	fileMock.readVal = "changed content"
	err = editor.Edit(s)
	assert.Nil(t, err)
	assert.Equal(t, 1, patcherCalled)

}


var mockRunner = func (a *types.Alias, app string, args []string, opts EditOptions) error {
	return nil
}

var patcher = func(called *int) SnippetPatcher {
	return func(req *types.PatchRequest) (*types.PatchResponse, error) {
		*called += 1
		return &types.PatchResponse{}, nil
	}
}