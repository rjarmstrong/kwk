package models

type SearchResult struct {
	Id          int64
	Username    string
	Key         string
	Runtime     string
	Extension   string
	Tags        []string
	Snip        string
	SnipVersion int64
	Description string
	Highlights  map[string]string
}

type SearchResponse struct {
	Results []*SearchResult
	Total   int64
	Took    int64
	Term    string
}
