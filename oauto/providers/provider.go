package providers

import (
	"net/http"
	"github.com/ibrt/go-oauto/oauto/providers/facebook"
	"sourcegraph.com/sourcegraph/go-selenium"
)

var Providers = map[string]Provider{
	"facebook": facebook.NewFacebook(),
}

type Provider interface {
	GetName() string
	HandleRedirect(r *http.Request) (string, error)
	Authenticate(driver selenium.WebDriver, appID, username, password, redirectURL string) error
}