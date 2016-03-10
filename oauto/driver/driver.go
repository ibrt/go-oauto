package driver

import (
	"github.com/ibrt/go-oauto/oauto/providers"
	"sourcegraph.com/sourcegraph/go-selenium"
	"github.com/go-errors/errors"
	"github.com/ibrt/go-oauto/oauto/config"
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
	tokenDivID = "token"
)

func PerformAuthentication(config *config.Config, provider providers.Provider, appID, username, password, redirectURL string) (string, error) {
	// Initialize Selenium.
	driver, err := selenium.NewRemote(seleniumCaps, config.SeleniumURL)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	defer driver.Quit()
	if err := driver.SetImplicitWaitTimeout(seleniumImplicitTimeoutMS); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Execute the provider's authentication flow.
	err = provider.Authenticate(driver, appID, username, password, redirectURL)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Extract the token from the redirect page.
	element, err := driver.FindElement(selenium.ById, tokenDivID)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	token, err := element.Text()
	if err != nil {
		return "", errors.Wrap(err, 0)
	}

	return token, nil
}