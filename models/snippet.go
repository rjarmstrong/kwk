package models

import (
	"time"
)

const (
	ProfileFullKey  = "profile.json"
	TokenHeaderName = "token"
	ROOT_POUCH = "root"
	SETUP_POUCH = "setup"
)

type SnippetList struct {
	Username string
	Items    []Snippet
	Total    int64
	Since    time.Time
	Size     int64
}

func NewSnippet(snippet string) *Snippet{
	return &Snippet{Snip:snippet, Alias:Alias{SnipName:SnipName{}}}
}

type Snippet struct {
	Id string

	FullName  string
	Alias

	Snip    string
	Version int64
	Media   string
	Tags    []string
	Created time.Time

	Description       string
	ClonedFromAlias   string
	ClonedFromVersion int64
	Private           bool
	CloneCount        int64
	RunCount          int64
	Role              SnipRole
}

type SnipRole int32

const (
	RoleStandard    SnipRole = 0
	RolePreferences SnipRole = 1
	RoleEnvironment SnipRole = 2
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
	All bool
	Username string
	Size int64
	Since int64
	Tags []string
}
