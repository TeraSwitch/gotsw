package gotsw

import (
	"context"
	"net/http"
)

const regionBasePath = "v1/region"

type RegionService interface {
	List(ctx context.Context) ([]Region, error)
}

type RegionsServiceHandler struct {
	client *Client
}

type regionsRoot struct {
	Regions []Region `json:"result"`
}

func (s *RegionsServiceHandler) List(ctx context.Context) ([]Region, error) {
	path := regionBasePath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(regionsRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Regions, nil
}
