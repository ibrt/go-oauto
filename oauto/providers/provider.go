package providers

import (
	"net/http"
	"github.com/ibrt/go-oauto/oauto/providers/facebook"
	"sourcegraph.com/sourcegraph/go-selenium"
	"github.com/ibrt/go-oauto/oauto/providers/google"
)

var Providers = map[string]Provider{
	"facebook": facebook.NewFacebook(),
	"google": google.NewGoogle(),
}

type Provider interface {
	GetName() string
	HandleRedirect(r *http.Request) (string, error)
	Authenticate(driver selenium.WebDriver, appID, appSecret, username, password, redirectURL string) (string, error)
}