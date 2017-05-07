package types

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Snippet struct {
	Alias
	Id                string
	Private           bool
	Snip              string
	Preview           string
	SnipVersion       int64
	Description       string
	Tags              []string
	ClonedFromAlias   string
	ClonedFromVersion int64
	ClonedFromId      string
	Latest            bool // THE LATEST VERSION OF THIS SNIPPET

	SnipChecksum  string
	SnipSignature string
	Encrypted     bool
	EncryptedAlgo string
	Role          SnipRole

	SnipCounts

	RunStatus     UseStatus
	RunStatusTime time.Time
	Dependencies  []string
	Apps          []string
	SupportedOs   []string
	Attribution   string
	Created       time.Time
}

func (st *Snippet) IsApp() bool {
	return len(st.Dependencies) > 0
}

func (st *Snippet) VerifyChecksum() bool {
	s := sha256.Sum256([]byte(st.Snip))
	actual := fmt.Sprintf("%x", s)
	return actual == st.SnipChecksum
}

type SnipCounts struct {
	Views  int64
	Runs   int64
	Clones int64
}
