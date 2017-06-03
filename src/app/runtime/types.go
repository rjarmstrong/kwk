package runtime

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
	"errors"
)

type SnippetPatcher func(req *types.PatchRequest) (*types.PatchResponse, error)
type SnippetGetter func(req *types.GetRequest) (*types.ListResponse, error)
type SnippetMaker func(req *types.CreateRequest) error
type UseLogger func(req *types.UseContext) (*types.LogUseResponse, error)
type DocGetter func() (string, error)

func newRuntimeAlias(username, name string, ext string, uniquePerMachine bool) *types.Alias {
	if uniquePerMachine {
		s, _ := os.Hostname()
		name = fmt.Sprintf("%s-%s", name, strings.ToLower(s))
	}
	return &types.Alias{
		Username: username,
		Pouch:    types.PouchSettings,
		Name:     name,
		Ext:      ext,
	}
}

func DefaultPrefs() *out.Prefs {
	return &out.Prefs{
		PrivateView:       true,
		AutoYes:           false,
		Covert:            false,
		RequireRunKeyword: false,
		Quiet:             false,
		SnippetThumbRows:  3,
		ExpandedThumbRows: 15,
		CommandTimeout:    60,
		RowSpaces:         true,
		RowLines:          true,
	}
}

func getPrefsAsString(prefs out.Prefs) string {
	ph := &PrefsFile{KwkPrefs: "v1", Options: prefs}
	b, err := yaml.Marshal(ph)
	if err != nil {
		out.LogErr(err)
	}
	return string(b)
}

func DefaultEnv() *yaml.MapSlice {
	env := &yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(defaultEnvString), env)
	if err != nil {
		log.Fatal(err)
	}
	return env
}

func GetSection(yml *yaml.MapSlice, name string) (*yaml.MapSlice, error) {
	rs, _ := getSubSection(yml, name)
	if rs == nil {
		return nil, errors.New(fmt.Sprintf("No %s section in given .yml", name))
	}
	return &rs, nil
}

func getSubSection(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("The yml config section '%s' is not valid please check it.", name)
		}
	}()
	f := func(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
		for _, v := range *yml {
			if v.Key == name {
				if slice, ok := v.Value.(yaml.MapSlice); ok {
					return slice, nil
				}
				if _, ok := v.Value.([]interface{}); ok {
					items := []string{}
					for _, v2 := range v.Value.([]interface{}) {
						items = append(items, v2.(string))
					}
					return nil, items
				}
				return nil, []string{v.Value.(string)}
			}
		}
		return nil, nil
	}
	sub, bottom := f(yml, name)
	if sub == nil && bottom == nil {
		return f(yml, "default")
	}
	return sub, bottom
}

type PrefsFile struct {
	KwkPrefs string
	Options  out.Prefs
}

/*
	Road-map:
	WipeTrail         bool  //deletes the history each time a command is run TODO: Security
	SessionTimeout    int64 // 0 = no timeout, TODO: Implement on api SECURITY
	AutoEncrypt       bool  //Encrypts all snippets when created. TODO: SECURITY
	RegulateUpdates   bool  //Updates based on the recommended schedule. If false get updates as soon as available.
	DisableRun        bool  //Completely disabled running scripts even if using -f TODO: Security
	Encrypt   	  bool // TODO: Security
	Decrypt   	  bool // TODO: Security
*/
