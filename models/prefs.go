package models

import (
	"gopkg.in/yaml.v2"
)

var prefs  *Preferences
var env   *yaml.MapSlice

func Prefs() *Preferences {
	return prefs
}

func Env() *yaml.MapSlice {
	return env
}

func SetPrefs(p *Preferences) {
	prefs = p
}

func SetEnv(e *yaml.MapSlice) {
	env = e
}

func DefaultPrefs() *Preferences {
	p := &Preferences{
		Global: false,
		AutoYes:false,
		Force:  false,
	}
	p.Debug = false
	p.Covert = false
	p.DisableRun = false
	p.WipeTrail = false
	p.SessionTimeout = 15
	p.ListAll = true
	p.ListLong = false
	p.RegulateUpdates = true
	return p
}

type PreferencesHolder struct {
	KwkPrefs    string
	Preferences PersistedPrefs
}

// Preferences is a container for file and flag preferences.
type Preferences struct {
	PersistedPrefs
	Global  bool // TODO: Implement in search api SEARCH
	//Quiet   bool // Only display fullNames (for multi deletes etc)
	AutoYes bool //TODO: Security
	Force   bool // TODO: Security Will squash warning messages e.g. when running third party snippets.
	Encrypt bool // TODO: Security
	Decrypt bool // TODO: Security
	ListLong bool
}

// PersistedPrefs are preferences which can be persistent.
type PersistedPrefs struct {
	Covert         bool  // Always opens browser in covert mode, when set to true flag should have no effect. TODO: Update env/darwin.yml
	ListAll        bool  //List all pouches including private. TODO: implement on api in search SEARCH.
	Debug          bool  // Displays more detailed output. TODO: Logging
	DisableRun     bool  //Completely disabled running scripts even if using -f TODO: Security
	WipeTrail      bool  //deletes the history each time a command is run TODO: Security
	SessionTimeout int64 // 0 = no timeout, TODO: Implement on api SECURITY
	AutoEncrypt    bool  //Encrypts all snippets when created. TODO: SECURITY
	RegulateUpdates bool //Updates based on the recommended schedule. If false get updates as soon as available.
}
