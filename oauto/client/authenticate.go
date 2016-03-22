package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ibrt/go-oauto/oauto/api"
	"io/ioutil"
	"net/http"
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

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Authenticate request failed with status %v: '%s'.", resp.StatusCode, respBytes)
	}

	authResp := &api.AuthenticateResponse{}
	if err := json.Unmarshal(respBytes, &authResp); err != nil {
		return nil, errors.WrapPrefix(err, string(respBytes), 0)
	}

	return authResp, nil
}
