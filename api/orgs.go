package api

import (
	"context"
	"errors"
	"net/url"

	"github.com/arigatomachine/cli/identity"
	"github.com/arigatomachine/cli/primitive"
)

// OrgsClient makes proxied requests to the registry's orgs endpoints
type OrgsClient struct {
	client *Client
}

// OrgResult is the payload returned for an org object
type OrgResult struct {
	ID      *identity.ID   `json:"id"`
	Version uint8          `json:"version"`
	Body    *primitive.Org `json:"body"`
}

type orgCreateRequest struct {
	Body struct {
		Name string `json:"name"`
	} `json:"body"`
}

// Create creates a new org with the given name. It returns the newly-created org.
func (o *OrgsClient) Create(ctx context.Context, name string) (*OrgResult, error) {
	org := orgCreateRequest{}
	org.Body.Name = name

	req, _, err := o.client.NewRequest("POST", "/orgs", nil, &org, true)
	if err != nil {
		return nil, err
	}

	res := OrgResult{}
	_, err = o.client.Do(ctx, req, &res, nil, nil)
	return &res, err
}

// GetByName retrieves an org by its named
func (o *OrgsClient) GetByName(ctx context.Context, name string) (*OrgResult, error) {
	v := &url.Values{}
	if name == "" {
		return nil, errors.New("invalid org name")
	}
	v.Set("name", name)

	req, _, err := o.client.NewRequest("GET", "/orgs", v, nil, true)
	if err != nil {
		return nil, err
	}

	orgs := []OrgResult{}
	_, err = o.client.Do(ctx, req, &orgs, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(orgs) < 1 {
		return nil, nil
	}

	return &orgs[0], nil
}

// List returns all organizations that the signed-in user has access to
func (o *OrgsClient) List(ctx context.Context) ([]OrgResult, error) {
	req, _, err := o.client.NewRequest("GET", "/orgs", nil, nil, true)
	if err != nil {
		return nil, err
	}

	orgs := []OrgResult{}
	_, err = o.client.Do(ctx, req, &orgs, nil, nil)
	return orgs, err
}

// RemoveMember removes a user from an org
func (o *OrgsClient) RemoveMember(ctx context.Context, orgID identity.ID,
	userID identity.ID) error {

	req, _, err := o.client.NewRequest("DELETE",
		"/orgs/"+orgID.String()+"/members/"+userID.String(), nil, nil, true)
	if err != nil {
		return err
	}

	_, err = o.client.Do(ctx, req, nil, nil, nil)
	return err
}
