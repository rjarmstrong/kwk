package models

import (
	"crypto/sha256"
	"fmt"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"strings"
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
	Created   int64

	Description       string
	ClonedFromAlias   string
	ClonedFromVersion int64
	Private           bool
	CloneCount        int64
	RunCount          int64
	ViewCount          int64
	Role              SnipRole

	RunStatus     UseStatus
	RunStatusTime int64
	Dependencies []string
	Attribution string

	CheckSum string
}

func (st *Snippet) VerifyChecksum() bool {
	s := sha256.Sum256([]byte(st.Snip))
	actual := fmt.Sprintf("%x", s)
	log.Debug("VERIFY CHECKSUM:%s = %s", actual, st.CheckSum)
	return actual == st.CheckSum
}

type SnipRole int32

type UseType int64
type UseStatus int64

const (
	UseStatusUnknown UseStatus = 0
	UseStatusSuccess UseStatus = 1
	UseStatusFail    UseStatus = 2

	UseTypeUnknown UseType = 0
	UseTypeView    UseType = 1
	UseTypeRun     UseType = 2
	UseTypeClone   UseType = 3

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
	IgnorePouches bool
	Category string
}

//TODO: Optimise this
func ScanVulnerabilities(snip string) error {
	if strings.Contains(snip, "rm -rf") || strings.Contains(snip, "rm ") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: Shell scripts cannot contain 'rm '.")
	}
	if strings.Contains(snip, ":(){") || strings.Contains(snip, "./$0|./$0&"){
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: Fork bomb detected.")
	}
	if strings.Contains(snip, "fork") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: the word 'fork' is not allowed in scripts.")
	}
	if strings.Contains(snip, "/dev/sd") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: '/dev/sd' is not allowed in scripts.")
	}
	if strings.Contains(snip, "/dev/null") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: '/dev/null' is not allowed in scripts.")
	}
	if strings.Contains(snip, "| sh") || strings.Contains(snip, "| bash") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: piping directly into terminal not allowed in scripts.")
	}
	if strings.Contains(snip, "nohup") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: 'nohup' command is not allowed.")
	}
	return nil
}
