package search

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type ISearch interface {
	Search(term string) (*models.SearchResponse, error)
}