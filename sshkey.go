package gotsw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// SSHKey represents an SSH key
type SSHKey struct {
	ID          int64   `json:"id"`
	Created     string  `json:"created"`
	Deleted     *string `json:"deleted,omitempty"`
	ObjectType  string  `json:"objectType"`  // The type of object (always "KEY")
	ProjectID   int64   `json:"projectId"`   // The project ID
	DisplayName string  `json:"displayName"` // The display name of the SSH key
	Key         string  `json:"key"`         // The SSH key content
}

// SshKeyResponse represents a response containing SSH key data
type SshKeyResponse Result[SSHKey]

// ListSshKeyResponse represents a response containing multiple SSH keys
type ListSshKeyResponse Result[[]SSHKey]

// ListSshKeys retrieves all SSH keys assigned to your project
func (c *Client) ListSshKeys(ctx context.Context) ([]SSHKey, error) {
	resp := &ListSshKeyResponse{}
	httpResp, err := c.Request(ctx, http.MethodGet, "SshKey", nil)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp.Result, nil
}

// GetSshKey retrieves a specific SSH key by ID
func (c *Client) GetSshKey(ctx context.Context, id int64) (SSHKey, error) {
	resp := &SshKeyResponse{}
	httpResp, err := c.Request(ctx, http.MethodGet, fmt.Sprintf("SshKey/%d", id), nil)
	if err != nil {
		return SSHKey{}, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return SSHKey{}, err
	}

	if !resp.Success {
		return SSHKey{}, errors.New(resp.Message)
	}

	return resp.Result, nil
}

type CreateSshKeyRequest struct {
	DisplayName string `json:"displayName"`
	ProjectID   int64  `json:"projectId"`
	Key         string `json:"key"`
}

// CreateSshKey creates a new SSH key
func (c *Client) CreateSshKey(ctx context.Context, projectID int64, key CreateSshKeyRequest) (SSHKey, error) {
	resp := &SshKeyResponse{}
	key.ProjectID = projectID
	httpResp, err := c.Request(ctx, http.MethodPost, "SshKey", key)
	if err != nil {
		return SSHKey{}, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return SSHKey{}, err
	}

	if !resp.Success {
		return SSHKey{}, errors.New(resp.Message)
	}

	return resp.Result, nil
}
