package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ibrt/go-oauto/oauto/config"
	"github.com/ibrt/go-oauto/oauto/driver"
	"github.com/ibrt/go-oauto/oauto/providers"
	"github.com/ibrt/go-oauto/oauto/redirect"
	"net/http"
)

type AuthenticateRequest struct {
	Provider  string `json:"provider"`
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}

func HandleAuthenticate(config *config.Config, r *http.Request, baseURL string) (interface{}, error) {
	// Parse request.
	req := &AuthenticateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	// Find provider.
	provider, ok := providers.Providers[req.Provider]
	if !ok {
		return nil, errors.Errorf("Unknown provider '%v'.", req.Provider)
	}

	// Perform authentication flow.
	redirectURL := fmt.Sprintf("%v%v", baseURL, redirect.MakePath(provider))
	token, err := driver.PerformAuthentication(config, provider, req.AppID, req.AppSecret, req.UserName, req.Password, redirectURL)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	// Return response.
	return &AuthenticateResponse{Token: token}, nil
}
