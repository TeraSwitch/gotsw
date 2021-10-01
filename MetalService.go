package gotsw

import (
	"context"
	"fmt"
	"net/http"
)

const metalBasePath = "v1/metal"

type MetalService interface {
	List(ctx context.Context, options *MetalListOptions) ([]Metal, error)
	Create(ctx context.Context, metalReq *MetalCreateRequest) (*Metal, error)
	Get(ctx context.Context, id uint64) (*Metal, error)

	PowerOff(ctx context.Context, id uint64) (*Metal, error)
	PowerOn(ctx context.Context, id uint64) (*Metal, error)
	Reset(ctx context.Context, id uint64) (*Metal, error)

	ListTemplates(ctx context.Context) ([]MetalTemplate, error)
	ListTiers(ctx context.Context) ([]MetalTier, error)
}

type MetalListOptions struct {
	RegionId string `url:"region,omitempty"`
	TierId   string `url:"tier,omitempty"`
	Tag      string `url:"tag,omitempty"`
	Limit    int    `url:"limit,omitempty"`
	Skip     int    `url:"skip,omitempty"`
}

type MetalServiceHandler struct {
	client *Client
}

type metalsRoot struct {
	Metals []Metal `json:"result"`
}

type metalRoot struct {
	Metal *Metal `json:"result"`
}

type metalTemplatesRoot struct {
	MetalTemplates []MetalTemplate `json:"result"`
}

type metalTiersRoot struct {
	MetalTiers []MetalTier `json:"result"`
}

type PowerCommandRequest struct {
	Command string `json:"command"`
}

type MetalCreateRequest struct {
	RegionId    string      `json:"regionId"`
	TierId      string      `json:"tierId"`
	DisplayName string      `json:"displayName"`
	ImageId     string      `json:"imageId"`
	Partitions  []Partition `json:"partitions"`
	RaidArrays  []RaidArray `json:"raidArrays"`
	SshKeyId    uint64      `json:"sshKeyId"`
	TemplateId  uint64      `json:"templateId"`
	Tags        []string    `json:"tags"`
}

func (s *MetalServiceHandler) List(ctx context.Context, options *MetalListOptions) ([]Metal, error) {
	path := metalBasePath
	path, err := addOptions(path, options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(metalsRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Metals, nil
}

func (s *MetalServiceHandler) PowerOff(ctx context.Context, id uint64) (*Metal, error) {
	path := fmt.Sprintf(metalBasePath+"/%d/powercommand", id)

	commandRequest := PowerCommandRequest{
		Command: "PowerOff",
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, commandRequest)
	if err != nil {
		return nil, err
	}

	root := new(metalRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Metal, err
}

func (s *MetalServiceHandler) PowerOn(ctx context.Context, id uint64) (*Metal, error) {
	path := fmt.Sprintf(metalBasePath+"/%d/powercommand", id)

	commandRequest := PowerCommandRequest{
		Command: "PowerOn",
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, commandRequest)
	if err != nil {
		return nil, err
	}

	root := new(metalRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Metal, err
}

func (s *MetalServiceHandler) Reset(ctx context.Context, id uint64) (*Metal, error) {
	path := fmt.Sprintf(metalBasePath+"/%d/powercommand", id)

	commandRequest := PowerCommandRequest{
		Command: "Reset",
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, commandRequest)
	if err != nil {
		return nil, err
	}

	root := new(metalRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Metal, err
}

func (s *MetalServiceHandler) Get(ctx context.Context, id uint64) (*Metal, error) {
	path := fmt.Sprintf(metalBasePath+"/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(metalRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Metal, nil
}

func (s *MetalServiceHandler) Create(ctx context.Context, createRequest *MetalCreateRequest) (*Metal, error) {

	req, err := s.client.NewRequest(ctx, http.MethodPost, metalBasePath, createRequest)
	if err != nil {
		return nil, err
	}

	root := new(metalRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Metal, err
}

func (s *MetalServiceHandler) ListTemplates(ctx context.Context) ([]MetalTemplate, error) {
	path := metalBasePath + "/templates"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(metalTemplatesRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.MetalTemplates, nil
}

func (s *MetalServiceHandler) ListTiers(ctx context.Context) ([]MetalTier, error) {
	path := metalBasePath + "/tiers"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(metalTiersRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.MetalTiers, nil
}
