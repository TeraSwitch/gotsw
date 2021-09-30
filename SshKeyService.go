package gotsw

import (
	"context"
	"net/http"
)

const sshKeyBasePath = "v1/sshkey"

type SshKeyService interface {
	List(ctx context.Context) ([]Region, error)
	Create(ctx context.Context, sshKeyReq *SshKeyCreateRequest) (*SshKey, error)
}

type SshKeyServiceHandler struct {
	client *Client
}

type sshKeysRoot struct {
	SshKeys []SshKey `json:"result"`
}

type sshKeyRoot struct {
	SshKey *SshKey `json:"result"`
}

type SshKeyCreateRequest struct {
	Name      string `json:"displayName"`
	PublicKey string `json:"key"`
}

func (s *SshKeyServiceHandler) List(ctx context.Context) ([]SshKey, error) {
	path := sshKeyBasePath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(sshKeysRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.SshKeys, nil
}

func (s *SshKeyServiceHandler) Create(ctx context.Context, createRequest *SshKeyCreateRequest) (*SshKey, error) {

	req, err := s.client.NewRequest(ctx, http.MethodPost, sshKeyBasePath, createRequest)
	if err != nil {
		return nil, err
	}

	root := new(sshKeyRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.SshKey, err
}
