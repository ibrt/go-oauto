package client

import (
	"github.com/ibrt/go-oauto/oauto/api"
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/go-errors/errors"
	"bytes"
)

func Authenticate(baseURL string, request *api.AuthenticateRequest) (*api.AuthenticateResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	resp, err := http.Post(fmt.Sprintf("%v/api/authenticate", baseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	authResp := &api.AuthenticateResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Authenticate request failed with status %v: '%+v'.", resp.StatusCode, authResp)
	}

	return authResp, nil
}
