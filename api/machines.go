package api

import (
	"context"

	"github.com/arigatomachine/cli/identity"
	"github.com/arigatomachine/cli/primitive"
)

// MachinesClient makes proxied requests to the registry's machines endpoints
type MachinesClient struct {
	client *Client
}

// MachinesResult is the payload returned for a keypair object
type MachinesResult struct {
	Machine *struct {
		ID   *identity.ID       `json:"id"`
		Body *primitive.Machine `json:"body"`
	}
	Memberships []*struct {
		ID   *identity.ID          `json:"id"`
		Body *primitive.Membership `json:"body"`
	}
	Tokens []*struct {
		ID   *identity.ID            `json:"id"`
		Body *primitive.MachineToken `json:"body"`
	}
}

type machinesCreateRequest struct {
	Name  string       `json:"name"`
	OrgID *identity.ID `json:"org_id"`
}

// Create a new machine in the given org
func (m *MachinesClient) Create(ctx context.Context, name string, orgID *identity.ID, output *ProgressFunc) error {
	mcr := machinesCreateRequest{
		Name:  name,
		OrgID: orgID,
	}

	req, reqID, err := m.client.NewRequest("POST", "/machines", nil, &mcr, false)
	if err != nil {
		return err
	}

	_, err = m.client.Do(ctx, req, nil, &reqID, output)
	return err
}
