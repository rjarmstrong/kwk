package search

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/rpc/src/searchRpc"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"google.golang.org/grpc"
)

type Search struct {
	client   searchRpc.SearchRpcClient
	headers  *rpc.Headers
}

func New(conn *grpc.ClientConn, s settings.ISettings, h *rpc.Headers) ISearch {
	return &Search{ client:searchRpc.NewSearchRpcClient(conn), headers:h}
}

func (s *Search) Search(term string) (*models.SearchResponse, error) {
	if res, err := s.client.Alpha(s.headers.GetContext(), &searchRpc.AlphaRequest{
		Term:term,
	}); err != nil {
		return nil, err
	} else {
		r := &models.SearchResponse{}
		r.Results = []*models.SearchResult{}
		r.Took = res.Took
		r.Total = res.Total
		for _, v := range res.Results {
			item := &models.SearchResult{}
			item.Key = v.Key
			item.Username = v.Username
			item.Uri = v.Uri
			item.Extension = v.Extension
			item.Highlights = v.Highlights
			r.Results = append(r.Results, item)
		}
		return r, nil
	}
}