package registry

import (
	"context"
	"log"
	"net/url"

	"github.com/arigatomachine/cli/envelope"
	"github.com/arigatomachine/cli/identity"
)

// Machines represents the `/machines` registry endpoint
type Machines struct {
	client *Client
}

// Post creates a new machine on the registry.
func (m *Machines) Post(ctx context.Context, pubKey, privKey,
	claim *envelope.Signed) (*envelope.Signed, *envelope.Signed, []envelope.Signed, error) {

	req, err := k.client.NewRequest("POST", "/keypairs", nil,
		ClaimedKeyPair{
			PublicKey:  pubKey,
			PrivateKey: privKey,
			Claims:     []envelope.Signed{*claim},
		})
	if err != nil {
		log.Printf("Error building http request: %s", err)
		return nil, nil, nil, err
	}

	resp := ClaimedKeyPair{}
	_, err = k.client.Do(ctx, req, &resp)
	if err != nil {
		log.Printf("Failed to create signing keypair: %s", err)
		return nil, nil, nil, err
	}

	return resp.PublicKey, resp.PrivateKey, resp.Claims, nil
}
