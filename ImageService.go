package gotsw

import (
	"context"
	"net/http"
)

const imageBasePath = "v1/image"

type ImageService interface {
	List(ctx context.Context) ([]Image, error)
}

type ImagesServiceHandler struct {
	client *Client
}

type imagesRoot struct {
	Images []Image `json:"result"`
}

func (s *ImagesServiceHandler) List(ctx context.Context) ([]Image, error) {
	path := imageBasePath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(imagesRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Images, nil
}