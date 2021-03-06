package api

import (
	"context"

	"github.com/arigatomachine/cli/apitypes"
)

// SessionClient provides access to the daemon's user session related endpoints,
// for logging in/out, and checking your session status.
type SessionClient struct {
	client *Client
}

// Get returns the status of the user's session.
func (s *SessionClient) Get(ctx context.Context) (*apitypes.SessionStatus, error) {
	req, _, err := s.client.NewRequest("GET", "/session", nil, nil, false)
	if err != nil {
		return nil, err
	}

	resp := &apitypes.SessionStatus{}
	_, err = s.client.Do(ctx, req, resp, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Login logs the user in using the provided email and passphrase
func (s *SessionClient) Login(ctx context.Context, email, passphrase string) error {
	login := apitypes.Login{
		Email:      email,
		Passphrase: passphrase,
	}
	req, _, err := s.client.NewRequest("POST", "/login", nil, &login, false)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil

}

// Logout logs the user out of their session
func (s *SessionClient) Logout(ctx context.Context) error {
	req, _, err := s.client.NewRequest("POST", "/logout", nil, nil, false)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
