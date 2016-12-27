package search

import "bitbucket.com/sharingmachine/kwkcli/models"

type ISearch interface {
	Search(term string) (*models.SearchResponse, error)
}
