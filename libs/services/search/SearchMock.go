package search

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type SearchMock struct {
}

func (s *SearchMock) Search(term string) (*models.SearchResponse, error) {
	return nil, nil
}