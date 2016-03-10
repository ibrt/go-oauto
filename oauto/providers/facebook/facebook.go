package facebook

import (
	"net/http"
	"sourcegraph.com/sourcegraph/go-selenium"
	"fmt"
	"strings"
	"github.com/go-errors/errors"
)

type Facebook struct {
	// Intentionally empty.
}

const (
	authURL = "https://www.facebook.com/dialog/oauth?client_id=%v&redirect_uri=%v"
	emailFieldName = "email"
	passwordFieldName = "pass"
	loginButtonName = "login"
	authorizeButtonName = "__CONFIRM__"
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

func (f *Facebook) Authenticate(driver selenium.WebDriver, appID, username, password, redirectURL string) error {
	// Load FB auth page.
	if err := driver.Get(fmt.Sprintf(authURL, appID, redirectURL)); err != nil {
		return errors.Wrap(err, 0)
	}

	// Fill e-mail and password fields, click "Login".
	element, err := driver.FindElement(selenium.ByName, emailFieldName)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	if err := element.SendKeys(username); err != nil {
		return errors.Wrap(err, 0)
	}
	element, err = driver.FindElement(selenium.ByName, passwordFieldName)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	if err := element.SendKeys(password); err != nil {
		return errors.Wrap(err, 0)
	}
	element, err = driver.FindElement(selenium.ByName, loginButtonName)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	if err := element.Click(); err != nil {
		return errors.Wrap(err, 0)
	}

	// If needed, click authorize the app. If the app is already authorized, just return.
	element, err = driver.FindElement(selenium.ByName, authorizeButtonName)
	if err != nil {
		if url, _ := driver.CurrentURL(); strings.HasPrefix(url, redirectURL) {
			return nil
		}
		return errors.Wrap(err, 0)
	}
	if err := element.Click(); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}