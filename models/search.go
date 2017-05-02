package models

import "bitbucket.com/sharingmachine/types"

type SearchResult struct {
	*types.Snippet
	Highlights  map[string]string
}

type SearchTermResponse struct {
	Results []*SearchResult
	Total   int64
	Took    int64
	Term    string
}
