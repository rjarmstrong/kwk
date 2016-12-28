package search

import "bitbucket.com/sharingmachine/kwkcli/models"

type TermMock struct {
}

func (s *TermMock) Execute(term string) (*models.SearchResponse, error) {
	return nil, nil
}
