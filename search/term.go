package search

import "bitbucket.com/sharingmachine/kwkcli/models"

// Term is a standard term search which depending on the
// implementation is analyzed in the kwk api.
type Term interface {
	Execute(term string) (*models.SearchTermResponse, error)
}
