package models

import (
	"time"
)

const (
	ProfileFullKey  = "profile.json"
	TokenHeaderName = "token"
	ROOT_POUCH      = ""
	SETTINGS_POUCH  = "settings"
)

func NewSnippet(snippet string) *Snippet {
	return &Snippet{Snip: snippet, Alias: Alias{SnipName: SnipName{}}}
}

type Snippet struct {
	Id string

	Alias

	Snip      string
	Signature string
	Version   int64
	Media     string
	Tags      []string
	Preview   string
	Created   time.Time

	Description       string
	ClonedFromAlias   string
	ClonedFromVersion int64
	Private           bool
	CloneCount        int64
	RunCount          int64
	Role              SnipRole

	RunStatus     RunStatus
	RunStatusTime int64
}

type SnipRole int32
type RunStatus int64

const (
	RunStatusUnknown RunStatus = 0
	RunStatusSuccess RunStatus = 1
	RunStatusFail    RunStatus = 2

	SnipRoleStandard    SnipRole = 0
	SnipRolePreferences SnipRole = 1
	SnipRoleEnvironment SnipRole = 2
)

type CreateSnippetResponse struct {
	Snippet   *Snippet
	TypeMatch *TypeMatch
}

type TypeMatch struct {
	Matches []Match
}

type Match struct {
	Score     int64
	Extension string
}

type ListParams struct {
	All      bool
	Pouch    string
	Username string
	Size     int64
	Since    int64
	Tags     []string
}
