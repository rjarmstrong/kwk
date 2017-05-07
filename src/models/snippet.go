package models

import (
	"github.com/kwk-super-snippets/types"
)

const (
	ProfileFullKey = "profile.json"
)

type CreateSnippetResponse struct {
	Snippet *types.Snippet
}

type ListParams struct {
	All           bool
	Pouch         string
	Username      string
	Size          int64
	Since         int64
	Tags          []string
	IgnorePouches bool
	Category      string
}
