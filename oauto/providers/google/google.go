package google

import (
	"fmt"
	"github.com/go-errors/errors"
	"net/http"
	"net/url"
	"sourcegraph.com/sourcegraph/go-selenium"
	"strings"
	"time"
)

type Google struct {
	// Intentionally empty.
}

const (
	authURL           = "https://accounts.google.com/o/oauth2/auth?client_id=%v&redirect_uri=%v&response_type=id_token&approval_prompt=force&scope=https://www.googleapis.com/auth/plus.me https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
	emailFieldName    = "Email"
	passwordFieldName = "Passwd"
	signInButtonName  = "signIn"
	authorizeButtonID = "submit_approve_access"
)

func NewGoogle() *Google {
	return &Google{}
}

func (g *Google) GetName() string {
	return "google"
}

func (g *Google) HandleRedirect(r *http.Request) (string, error) {
	// In this case we get the token from the URL fragment (which is not passed to the server).
	return "", nil
}

func (g *Google) Authenticate(webDriver selenium.WebDriver, appID, appSecret, username, password, redirectURL string) (string, error) {
	// Load Google auth page.
	if err := webDriver.Get(fmt.Sprintf(authURL, appID, redirectURL)); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Fill e-mail field, click "Next".
	element, err := webDriver.FindElement(selenium.ByName, emailFieldName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.SendKeys(username); err != nil {
		return "", errors.Wrap(err, 0)
	}
	element, err = webDriver.FindElement(selenium.ByName, signInButtonName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.Click(); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Fill password field, click "Sign in".
	element, err = webDriver.FindElement(selenium.ByName, passwordFieldName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.SendKeys(password); err != nil {
		return "", errors.Wrap(err, 0)
	}
	element, err = webDriver.FindElement(selenium.ById, signInButtonName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.Click(); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// If needed, go through the extra security flow.
	element, err = webDriver.FindElement(selenium.ById, "MapChallenge")
	if err == nil {
		if err := element.Click(); err != nil {
			return "", errors.Wrap(err, 0)
		}
		element, err = webDriver.FindElement(selenium.ById, "address")
		if err := element.SendKeys("San Francisco"); err != nil {
			return "", errors.Wrap(err, 0)
		}
		element, err = webDriver.FindElement(selenium.ById, "submitChallenge")
		if err != nil {
			return "", errors.Wrap(err, 0)
		}
		if err := element.Click(); err != nil {
			return "", errors.Wrap(err, 0)
		}
	}

	// If needed, click authorize the app. If the app is already authorized, just continue.
	element, err = webDriver.FindElement(selenium.ById, authorizeButtonID)
	if err == nil {

		url, _ := webDriver.CurrentURL()
		fmt.Printf("%s\n", url)
		src, _ := webDriver.PageSource()
		fmt.Printf("%s\n", src)

		// Wait for the button to become clickable.
		time.Sleep(2 * time.Second)
		if err := element.Click(); err != nil {
			return "", errors.Wrap(err, 0)
		}
	} else {

		url, _ := webDriver.CurrentURL()
		fmt.Printf("%s\n", url)
		src, _ := webDriver.PageSource()
		fmt.Printf("%s\n", src)

		if url, _ := webDriver.CurrentURL(); !strings.HasPrefix(url, redirectURL) {
			return "", errors.Wrap(err, 0)
		}
	}

	// Wait until redirect is complete.
	currentURL, err := waitForCurrentURLPrefix(webDriver, redirectURL)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Parse redirect URL and extract id_token from fragment.
	parsedURL, err := url.Parse(currentURL)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	values, err := url.ParseQuery(parsedURL.Fragment)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if values.Get("id_token") == "" {
		return "", errors.New("Missing 'id_token' key in redirect URL fragment.")
	}

	return values.Get("id_token"), nil
}

func waitForCurrentURLPrefix(webDriver selenium.WebDriver, prefix string) (string, error) {
	for i := 0; i < 10; i++ {
		currentURL, err := webDriver.CurrentURL()
		if err != nil {
			return "", errors.Wrap(err, 0)
		}
		if strings.HasPrefix(currentURL, prefix) {
			return currentURL, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return "", fmt.Errorf("Current URL never had prefix '%v'.", prefix)
}
