package providers

import (
	"github.com/ibrt/go-oauto/oauto/providers/digits"
	"github.com/ibrt/go-oauto/oauto/providers/facebook"
	"github.com/ibrt/go-oauto/oauto/providers/google"
	"net/http"
	"sourcegraph.com/sourcegraph/go-selenium"
)

var Providers = map[string]Provider{
	"facebook": facebook.NewFacebook(),
	"google":   google.NewGoogle(),
	"digits":   digits.NewDigits(),
}

type Provider interface {
	GetName() string
	HandleRedirect(r *http.Request) (string, error)
	Authenticate(driver selenium.WebDriver, appID, appSecret, username, password, redirectURL string) (string, error)
}
