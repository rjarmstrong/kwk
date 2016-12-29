package search

import "bitbucket.com/sharingmachine/kwkcli/models"

type TermMock struct {
	ReturnForExecute *models.SearchTermResponse
}

func (s *TermMock) Execute(term string) (*models.SearchTermResponse, error) {
	if s.ReturnForExecute != nil {
		return s.ReturnForExecute, nil
	}
	return nil, nil
}
