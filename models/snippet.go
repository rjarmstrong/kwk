package models

import "time"

const (
	ProfileFullKey  = "profile.json"
	TokenHeaderName = "token"
)

type Alias struct {
	Username string `json:"username" schema:"username"`
	FullKey  string `json:"fullKey" schema:"fullKey"`
}

type SnippetList struct {
	Username string
	Items    []Snippet
	Total    int64
	Since    time.Time
	Size     int64
}

type Snippet struct {
	Id string

	FullName  string
	Username  string
	Name      string
	Extension string

	Snip    string
	Version int64
	Runtime string
	Media   string
	Tags    []string
	Created time.Time

	Description        string
	ClonedFromFullName string
	ClonedFromVersion  int64
	Private            bool
	CloneCount         int64
	RunCount           int64
}

type CreateSnippetRequest struct {
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
