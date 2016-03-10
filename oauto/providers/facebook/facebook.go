package facebook

import (
	"net/http"
	"sourcegraph.com/sourcegraph/go-selenium"
	"fmt"
	"strings"
	"github.com/go-errors/errors"
	"encoding/json"
)

type Facebook struct {
	// Intentionally empty.
}

const (
	authURL = "https://www.facebook.com/dialog/oauth?client_id=%v&redirect_uri=%v?scope=public_profile,email"
	exchangeURL = "https://graph.facebook.com/v2.3/oauth/access_token?client_id=%v&redirect_uri=%v&client_secret=%v&code=%v"
	emailFieldName = "email"
	passwordFieldName = "pass"
	loginButtonName = "login"
	authorizeButtonName = "__CONFIRM__"
	tokenDivID = "token"
)

func NewFacebook() *Facebook {
	return &Facebook{}
}

func (f *Facebook) GetName() string {
	return "facebook"
}

func (f *Facebook) HandleRedirect(r *http.Request) (string, error) {
	if token := r.URL.Query().Get("code"); len(token) > 0 {
		return token, nil
	} else {
		return "", errors.Errorf("Missing 'code' query string parameter in request '%s'.", r.URL)
	}
}

func (f *Facebook) Authenticate(driver selenium.WebDriver, appID, appSecret, username, password, redirectURL string) (string, error) {
	// Load FB auth page.
	if err := driver.Get(fmt.Sprintf(authURL, appID, redirectURL)); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Fill e-mail and password fields, click "Login".
	element, err := driver.FindElement(selenium.ByName, emailFieldName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.SendKeys(username); err != nil {
		return "", errors.Wrap(err, 0)
	}
	element, err = driver.FindElement(selenium.ByName, passwordFieldName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.SendKeys(password); err != nil {
		return "", errors.Wrap(err, 0)
	}
	element, err = driver.FindElement(selenium.ByName, loginButtonName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.Click(); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// If needed, click authorize the app. If the app is already authorized, just continue.
	element, err = driver.FindElement(selenium.ByName, authorizeButtonName)
	if err == nil {
		if err := element.Click(); err != nil {
			return "", errors.Wrap(err, 0)
		}
	} else {
		if url, _ := driver.CurrentURL(); !strings.HasPrefix(url, redirectURL) {
			return "", errors.Wrap(err, 0)
		}
	}

	// Extract code token from the redirect page.
	element, err = driver.FindElement(selenium.ById, tokenDivID)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	code, err := element.Text()
	if err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Exchange the code for a OAuth token using the secret app id.
	resp, err := http.Get(fmt.Sprintf(exchangeURL, appID, redirectURL, appSecret, code))
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("Token exchange request failed with status %v.", resp.StatusCode)
	}
	exchangeResp := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return "", errors.Wrap(err, 0)
	}

	return exchangeResp["access_token"].(string), nil
}