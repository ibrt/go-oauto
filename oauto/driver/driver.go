package driver

import (
	"github.com/go-errors/errors"
	"github.com/ibrt/go-oauto/oauto/config"
	"github.com/ibrt/go-oauto/oauto/providers"
	"sourcegraph.com/sourcegraph/go-selenium"
)

var seleniumCaps = selenium.Capabilities(map[string]interface{}{
	"browserName": "chrome",
	"chromeOptions": map[string]interface{}{
		"prefs": map[string]interface{}{
			"profile.default_content_setting_values.notifications": 2,
		},
	},
})

const (
	seleniumImplicitTimeoutMS = 10000
)

func PerformAuthentication(config *config.Config, provider providers.Provider, appID, appSecret, username, password, redirectURL string) (string, error) {
	// Initialize Selenium.
	webDriver, err := selenium.NewRemote(seleniumCaps, config.SeleniumURL)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	defer webDriver.Quit()
	if err := webDriver.SetImplicitWaitTimeout(seleniumImplicitTimeoutMS); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Execute the provider's authentication flow.
	token, err := provider.Authenticate(webDriver, appID, appSecret, username, password, redirectURL)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}

	return token, nil
}
