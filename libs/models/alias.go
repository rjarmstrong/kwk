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
	Items []Alias `json:"items"`
	Total int32 `json:"total"`
	Page  int32 `json:"page"`
	Size  int32 `json:"size"`
}

type Alias struct {
	Id        int64 `json:"id"`

	FullKey   string `json:"fullKey"`
	Username  string `json:"username"`
	Key       string `json:"key"`
	Extension string `json:"extension"`

	Uri       string `json:"uri"`
	Version   int32 `json:"version"`
	Runtime   string `json:"runtime"`
	Media     string `json:"media"`
	Tags      []string `json:"tags"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type CreateAlias struct {
	Alias *Alias
	TypeMatch *TypeMatch
}

type TypeMatch struct {
	Matches []Match
}

type Match struct {
	Score     int32
	Media     string
	Runtime   string
	Extension string
}