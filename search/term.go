package search

import "bitbucket.com/sharingmachine/kwkcli/models"

type Term interface {
	Execute(term string) (*models.SearchResponse, error)
}
