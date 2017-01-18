package config

type FlagAndFile struct {
	Verbose bool // displays more detailed output
	Covert bool // opens browser in covert always
	HidePrivate bool //hides private snippets when searching or listing in cli
}

// Hydrate from disk
type FileConfig struct {
	FlagAndFile
	WarnExec3rdPartySnippets bool
	DisableRun bool
	AutoCrlf bool //Automatically change crlf files to lf (and vice versa?)
	WipeTrail bool //deletes the history each time a command is run
}

// Hydrate from disk but override if given
type Flags struct {
	FileConfig
	FlagAndFile

	Global  bool
	Quiet bool // Only display fullNames (for multi deletes etc)
	AutoYes bool
	Force bool
}
