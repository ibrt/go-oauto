package oauto_test

import (
	"testing"
	"github.com/kelseyhightower/envconfig"
	"os"
	"github.com/ibrt/go-oauto/oauto/api"
	"github.com/ibrt/go-oauto/oauto/client"
	"github.com/go-errors/errors"
)

type TestConfig struct {
	BaseURL string `envconfig:"BASE_URL" required:"true"`
	FacebookAppID string `envconfig:"FACEBOOK_APP_ID" required:"true"`
	FacebookAppSecret string `envconfig:"FACEBOOK_APP_SECRET" required:"true"`
	FacebookUserName string `envconfig:"FACEBOOK_USER_NAME" required:"true"`
	FacebookPassword string `envconfig:"FACEBOOK_PASSWORD" required:"true"`
	GoogleAppID string `envconfig:"GOOGLE_APP_ID" required:"true"`
	GoogleAppSecret string `envconfig:"GOOGLE_APP_SECRET" required:"true"`
	GoogleUserName string `envconfig:"GOOGLE_USER_NAME" required:"true"`
	GooglePassword string `envconfig:"GOOGLE_PASSWORD" required:"true"`
}

var testConfig = &TestConfig{}

func TestMain(m *testing.M) {
	envconfig.MustProcess("OAUTO_TEST", testConfig)
	os.Exit(m.Run())
}

func TestFacebook(t *testing.T) {
	resp, err := client.Authenticate(
		testConfig.BaseURL,
		&api.AuthenticateRequest{
			Provider: "facebook",
			AppID: testConfig.FacebookAppID,
			AppSecret: testConfig.FacebookAppSecret,
			UserName: testConfig.FacebookUserName,
			Password: testConfig.FacebookPassword,
		})
	if err != nil {
		t.Fatal(err.(*errors.Error).ErrorStack())
	}
	if len(resp.Token) == 0 {
		t.Fatal("Missing token in Facebook authentication response.")
	}
}

/*func TestGoogle(t *testing.T) {
	resp, err := client.Authenticate(
		testConfig.BaseURL,
		&api.AuthenticateRequest{
			Provider: "google",
			AppID: testConfig.GoogleAppID,
			AppSecret: testConfig.GoogleAppSecret,
			UserName: testConfig.GoogleUserName,
			Password: testConfig.GooglePassword,
		})
	if err != nil {
		t.Fatal(err.(*errors.Error).ErrorStack())
	}
	if len(resp.Token) == 0 {
		t.Fatal("Missing token in Google authentication response.")
	}
}*/