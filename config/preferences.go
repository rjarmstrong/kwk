package config

func NewPreferences() *Preferences {
	p := &Preferences{
		Global: false,
		Quiet:  false,
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
	Global  bool
	Quiet   bool // Only display fullNames (for multi deletes etc)
	AutoYes bool
	Force   bool // Will squash warning messages e.g. when running third party snippets.
	Encrypt bool
	Decrypt bool
}

// PersistedPrefs are preferences which can be persistent.
type PersistedPrefs struct {
	Covert      bool // Always opens browser in covert mode, when set to true flag should have no effect.
	Verbose     bool // Displays more detailed output
	HidePrivate bool //hides private snippet text when searching or listing in cli
	DisableRun  bool //Completely disabled running scripts even if using -f
	WipeTrail      bool  //deletes the history each time a command is run
	SessionTimeout int64 // 0 = no timeout, TODO: Implement on api
	AutoEncrypt    bool  //Encrypts all snippets when created.
}
