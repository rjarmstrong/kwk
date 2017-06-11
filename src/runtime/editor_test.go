package runtime

import (
	"crypto/sha256"
	"fmt"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk-types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

var content = "some content"

func Test_Editor(t *testing.T) {
	out.DebugEnabled = true
	fileMock := &snippetReadWriter{}
	//username := "test-man"
	env := &yaml.MapSlice{}
	prefs := &out.Prefs{}
	yaml.Unmarshal([]byte(testEnvString), env)
	var patcherCalled int
	fileMock.readVal = content
	editor := NewEditor(env, prefs, patcher(&patcherCalled), fileMock).(*editor)
	editor.gui = mockRunner
	s := types.NewBlankSnippet()
	s.Content = content
	shaa := sha256.Sum256([]byte(content))
	s.Checksum = fmt.Sprintf("%x", shaa)

	t.Log("If file unchanged should Not patch file.")
	onchange := func(s types.Snippet) {}
	err := editor.Invoke(s, onchange)
	assert.Nil(t, err)
	changes, err := editor.Close(s)
	assert.Nil(t, err)
	assert.Equal(t, 0, patcherCalled)
	assert.Equal(t, uint(0), changes)
	assert.Equal(t, 1, fileMock.rmDirCalled)

	t.Log("If file changed Should patch file.")
	fileMock.readVal = "changed content"
	fileMock.rmDirCalled = 0
	err = editor.Invoke(s, onchange)
	assert.Nil(t, err)
	changes, err = editor.Close(s)
	assert.Nil(t, err)
	assert.Equal(t, 1, patcherCalled)
	assert.Equal(t, uint(1), changes)
	assert.Equal(t, 1, fileMock.rmDirCalled)
}

var mockRunner = func(a *types.Alias, app string, args []string, opts EditOptions) error {
	return nil
}

var patcher = func(called *int) SnippetPatcher {
	return func(req *types.PatchRequest) (*types.PatchResponse, error) {
		*called += 1
		snippet := types.NewBlankSnippet()
		snippet.Content = content
		return &types.PatchResponse{Snippet: snippet}, nil
	}
}
