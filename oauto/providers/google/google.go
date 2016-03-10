package google

import (
	"net/http"
	"sourcegraph.com/sourcegraph/go-selenium"
	"fmt"
	"github.com/go-errors/errors"
	"time"
	"strings"
	"encoding/json"
	"bytes"
	"net/url"
)

type Google struct {
	// Intentionally empty.
}

const (
	authURL = "https://accounts.google.com/o/oauth2/auth?client_id=%v&redirect_uri=%v&response_type=code&approval_prompt=force&scope=email"
	exchangeURL = "https://accounts.google.com/o/oauth2/token"
	emailFieldName = "Email"
	passwordFieldName = "Passwd"
	loginButtonName="signIn"
	authorizeButtonID = "submit_approve_access"
	tokenDivID = "token"
)

func NewGoogle() *Google {
	return &Google{}
}

func (g *Google) GetName() string {
	return "google"
}

func (g *Google) HandleRedirect(r *http.Request) (string, error) {
	if token := r.URL.Query().Get("code"); len(token) > 0 {
		return token, nil
	} else {
		return "", errors.Errorf("Missing 'code' query string parameter in request '%s'.", r.URL)
	}
}

func (g *Google) Authenticate(driver selenium.WebDriver, appID, appSecret, username, password, redirectURL string) (string, error) {
	// Load Google auth page.
	if err := driver.Get(fmt.Sprintf(authURL, appID, redirectURL)); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Fill e-mail and password fields, click "Sign in".
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
	element, err = driver.FindElement(selenium.ById, authorizeButtonID)
	if err == nil {
		if err := element.Click(); err != nil {
			return "", errors.Wrap(err, 0)
		}
		// Wait for the button to become clickable.
		time.Sleep(10 * time.Second)
		if err := element.Click(); err != nil {
			return "", errors.Wrap(err, 0)
		}
	} else {
		if url, _ := driver.CurrentURL(); !strings.HasPrefix(url, redirectURL) {
			return "", errors.Wrap(err, 0)
		}
	}

	// Extract the code from the redirect page.
	element, err = driver.FindElement(selenium.ById, tokenDivID)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	code, err := element.Text()
	if err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Build the exchange request body.
	exchangeBody := url.Values{}
	exchangeBody.Set("client_id", appID)
	exchangeBody.Set("client_secret", appSecret)
	exchangeBody.Set("redirect_uri", redirectURL)
	exchangeBody.Set("code", code)
	exchangeBody.Set("grant_type", "authorization_code")

	// Exchange the code for a OAuth token using the secret app id.
	resp, err := http.Post(exchangeURL, "application/x-www-form-urlencoded", bytes.NewBufferString(exchangeBody.Encode()))
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if resp.StatusCode != http.StatusOK {
		//return "", errors.Errorf("Token exchange request failed with status %v.", resp.StatusCode)
	}
	exchangeResp := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return "", errors.Wrap(err, 0)
	}

	return exchangeResp["access_token"].(string), nil
}