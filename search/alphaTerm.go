package search

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/rpc/src/searchRpc"
	"google.golang.org/grpc"
)

type AlphaTerm struct {
	client  searchRpc.SearchRpcClient
	headers *rpc.Headers
}

func NewAlphaTerm(conn *grpc.ClientConn, s config.Persister, h *rpc.Headers) Term {
	return &AlphaTerm{client: searchRpc.NewSearchRpcClient(conn), headers: h}
}

func (s *AlphaTerm) Execute(term string) (*models.SearchTermResponse, error) {
	if res, err := s.client.Alpha(s.headers.GetContext(), &searchRpc.AlphaRequest{
		Term: term,
	}); err != nil {
		return nil, err
	} else {
		r := &models.SearchTermResponse{}
		r.Results = []*models.SearchResult{}
		r.Took = res.Took
		r.Total = res.Total
		for _, v := range res.Results {
			item := &models.SearchResult{}
			item.Name = v.Name
			item.FullName = v.FullName
			item.Username = v.Username
			item.Description = v.Description
			item.SnipVersion = v.SnipVersion
			item.Snip = v.Snip
			item.Extension = v.Extension
			item.Highlights = v.Highlights
			item.RunCount = v.RunCount
			item.CloneCount = v.CloneCount
			item.Created = v.Created
			item.Role = models.SnipRole(v.Role)
			item.Private = v.Private
			item.Tags = v.Tags
			r.Results = append(r.Results, item)
		}
		return r, nil
	}
}
