package models

type SearchResult struct {
	*Snippet
	Highlights  map[string]string
}

type SearchTermResponse struct {
	Results []*SearchResult
	Total   int64
	Took    int64
	Term    string
}
