package models

import "time"

const (
	ProfileFullKey  = "profile.json"
	TokenHeaderName = "token"
)

type KwkKey struct {
	Username string `json:"username" schema:"username"`
	FullKey  string `json:"fullKey" schema:"fullKey"`
}

type SnippetList struct {
	Username string  `json:"username"`
	Items    []Snippet `json:"items"`
	Total    int64   `json:"total"`
	Since    time.Time
	Size     int64
}

type Snippet struct {
	Id string `json:"id"`

	FullKey   string `json:"fullKey"`
	Username  string `json:"username"`
	Key       string `json:"key"`
	Extension string `json:"extension"`

	Snip    string    `json:"uri"`
	Version int64     `json:"version"`
	Runtime string    `json:"runtime"`
	Media   string    `json:"media"`
	Tags    []string  `json:"tags"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Description       string
	ForkedFromFullKey string
	ForkedFromVersion int64
	Private           bool
	CloneCount        int64
	RunCount          int64
}

type CreateSnippet struct {
	Alias     *Snippet
	TypeMatch *TypeMatch
}

type TypeMatch struct {
	Matches []Match
}

type Match struct {
	Score     int64
	Media     string
	Runtime   string
	Extension string
}
