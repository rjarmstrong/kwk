package models

import "time"

const (
	ProfileFullKey = "profile.json"
	TokenHeaderName = "token"
)

type KwkKey struct {
	Username string `json:"username" schema:"username"`
	FullKey string `json:"fullKey" schema:"fullKey"`
}

type AliasList struct {
	Username string `json:"username"`
	Items []Alias `json:"items"`
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
}

type Alias struct {
	Id        int64 `json:"id"`

	FullKey   string `json:"fullKey"`
	Username  string `json:"username"`
	Key       string `json:"key"`
	Extension string `json:"extension"`

	Uri       string `json:"uri"`
	Version   int64 `json:"version"`
	Runtime   string `json:"runtime"`
	Media     string `json:"media"`
	Tags      []string `json:"tags"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`

	Description       string
	ForkedFromFullKey string
	ForkedFromVersion int64
	Private           bool
	ForkCount         int64
	RunCount          int64
}

type CreateAlias struct {
	Alias *Alias
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