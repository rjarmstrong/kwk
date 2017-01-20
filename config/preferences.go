package config

func DefaultPrefs() *Preferences {
	p := &Preferences{
		Global: false,
		AutoYes:false,
		Force:  false,
	}
	p.Verbose = false
	p.Covert = false
	p.HidePrivate = true
	p.DisableRun = false
	p.WipeTrail = false
	p.SessionTimeout = 15
	return p
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
}

// PersistedPrefs are preferences which can be persistent.
type PersistedPrefs struct {
	Covert      bool // Always opens browser in covert mode, when set to true flag should have no effect. TODO: Update env/darwin.yml
	Verbose     bool // Displays more detailed output. TODO: Logging
	HidePrivate bool //hides content of private snippet text when searching or listing in cli TODO: Security
	DisableRun  bool //Completely disabled running scripts even if using -f TODO: Security
	WipeTrail      bool  //deletes the history each time a command is run TODO: Security
	SessionTimeout int64 // 0 = no timeout, TODO: Implement on api SECURITY
	AutoEncrypt    bool  //Encrypts all snippets when created. TODO: SECURITY
	ListAll bool //List private files all the time. TODO: implement on api in search SEARCH.
}
