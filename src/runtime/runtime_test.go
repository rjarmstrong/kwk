package runtime

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/types"
	"github.com/stretchr/testify/assert"
	"github.com/kwk-super-snippets/types/errs"
)

func Test_Runtime(t *testing.T) {
	out.DebugEnabled = true
	fileMock := &snippetReadWriter{}
	username := "test-man"
	eh := &ehMock{}

	env := &yaml.MapSlice{}
	prefs := &out.Prefs{}
	testEnv := &yaml.MapSlice{}
	yaml.Unmarshal([]byte(testEnvString), testEnv)

	t.Log("Not logged in")
	Configure(env, prefs, "", getter(nil), maker(nil), fileMock, eh)
	assert.EqualValues(t, DefaultPrefs(), prefs)
	assert.EqualValues(t, DefaultEnv(), env)

	t.Log("Logged in and found in local cache")
	fileMock.readVal = getPrefsAsString(out.Prefs{ExpandedThumbRows: 7})
	Configure(env, prefs, username, getter(nil), maker(nil),  fileMock, eh)
	assert.EqualValues(t, 7, prefs.ExpandedThumbRows)

	fileMock.readVal = testEnvString
	Configure(env, prefs, username, getter(nil), maker(nil),  fileMock, eh)
	assert.EqualValues(t, testEnv, env)

	t.Log("Not found in local cache but found remote")
	fileMock.readVal = ""
	var makerCalled int
	var getterCalled int
	Configure(env, prefs, username, getter(&getterCalled), maker(&makerCalled),  fileMock, eh)
	assert.EqualValues(t, 16, prefs.SnippetThumbRows)
	assert.EqualValues(t, testEnv, env)
	assert.EqualValues(t, 2, getterCalled)
	assert.EqualValues(t, 0, makerCalled)


	t.Log("Remote Not found, creates snippet with default content and writes to local cache")
	fileMock.readVal = ""
	fileMock.writeCalledWith = []string{}
	Configure(env, prefs, username, nilGetter, maker(&makerCalled),  fileMock, eh)
	assert.Equal(t, 2, makerCalled)
	makerCalled = 0
	// Write called twice, we are only checking the most recent call here
	assert.EqualValues(t, 2, len(fileMock.writeCalledWith))
}

var nilGetter = func(req *types.GetRequest) (*types.ListResponse, error) {
	return nil, errs.NotFound
}

var getter = func (called *int) SnippetGetter {
	return func(req *types.GetRequest) (*types.ListResponse, error) {
		*called += 1
		if req.Alias.Name == "prefs" {
			return &types.ListResponse{Items: []*types.Snippet{{Content: getPrefsAsString(out.Prefs{SnippetThumbRows: 16})}}}, nil
		}
		return &types.ListResponse{Items: []*types.Snippet{{Content: testEnvString}}}, nil
	}
}

var maker = func(called *int) SnippetMaker {
	return func(req *types.CreateRequest) error {
		*called += 1
		return nil
	}
}

type ehMock struct {
}

func (*ehMock) Handle(err error) {
	panic("implement me")
}

type snippetReadWriter struct {
	readVal         string
	writeCalledWith []string
	rmDirCalled int
}

func (sm *snippetReadWriter) Write(uri string, content string) (string, error) {
	sm.writeCalledWith = append(sm.writeCalledWith, uri)
	return "", nil
}

func (sm *snippetReadWriter) Read(uri string) (string, error) {
	if sm.readVal == "" {
		return "", errs.NotFound
	}
	return sm.readVal, nil
}

func (sm *snippetReadWriter) RmDir(uri string) error {
	sm.rmDirCalled++
	return nil
}

const testEnvString = `kwkenv: "1"
editors:
  default: ["gimp"]
apps:
  gimp: ["open", "-a", "gimp", "$DIR"]
runners:
  sh: ["/bin/zsh", "-c", "$SNIP"]`
