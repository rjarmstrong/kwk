package search

import "bitbucket.com/sharingmachine/kwkcli/models"

type TermMock struct {
}

func (s *TermMock) Execute(term string) (*models.SearchTermResponse, error) {
	return nil, nil
}
