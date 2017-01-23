package models

type SearchResult struct {
	Id          int64
	Username    string
	Name        string
	FullName    string
	Extension   string
	Tags        []string
	Snip        string
	SnipVersion int64
	Description string
	RunCount    int64
	Private     bool
	CloneCount  int64
	Role	    SnipRole
	Created	    int64
	Highlights  map[string]string
}

type SearchTermResponse struct {
	Results []*SearchResult
	Total   int64
	Took    int64
	Term    string
}
