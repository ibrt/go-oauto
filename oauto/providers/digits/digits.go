package digits

import (
	"fmt"
	"github.com/go-errors/errors"
	"net/http"
	"sourcegraph.com/sourcegraph/go-selenium"
	"time"
)

type Digits struct {
	// Intentionally empty.
}

const (
	authURL              = "https://www.digits.com/login?consumer_key=%v&host=%v"
	frameName            = "digits"
	phoneNumberFieldName = "x_auth_phone_number"
	sendButtonClass      = "submit-phone"
)

func NewDigits() *Digits {
	return &Digits{}
}

func (g *Digits) GetName() string {
	return "digits"
}

func (g *Digits) HandleRedirect(r *http.Request) (string, error) {
	if token := r.URL.Query().Get("code"); len(token) > 0 {
		return token, nil
	} else {
		return "", errors.Errorf("Missing 'code' query string parameter in request '%s'.", r.URL)
	}
}

func (g *Digits) Authenticate(driver selenium.WebDriver, appID, appSecret, username, password, redirectURL string) (string, error) {
	// Load Digits auth page.
	if err := driver.Get(fmt.Sprintf(authURL, appID, redirectURL)); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Locate form iframe.
	if err := driver.SwitchFrame(frameName); err != nil {
		return "", errors.Wrap(err, 0)
	}

	// Fill phone number field, click "Send confirmation code".
	element, err := driver.FindElement(selenium.ByName, phoneNumberFieldName)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.SendKeys(username); err != nil {
		return "", errors.Wrap(err, 0)
	}
	element, err = driver.FindElement(selenium.ByClassName, sendButtonClass)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	if err := element.Click(); err != nil {
		return "", errors.Wrap(err, 0)
	}

	time.Sleep(1 * time.Minute)

	// Extract code token from the redirect page.
	/*element, err := driver.FindElement(selenium.ById, tokenDivID)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	token, err := element.Text()
	if err != nil {
		return "", errors.Wrap(err, 0)
	}*/

	return "", nil
}
