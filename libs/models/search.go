package models

type SearchResult struct {
	Id int64
	Username string
	Key string
	Runtime string
	Extension string
	Tags []string
	Uri string
	Version int64
	CreatedUTC string
	UpdatedUTC string
}

type SearchResponse struct {
	Results []*SearchResult
	Total   int64
	Took    int64
}
