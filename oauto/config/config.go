package config

type Config struct {
	SeleniumURL string `envconfig:"SELENIUM_URL" required:"http://localhost:4444/wd/hub"`
	RedirectHost string `envconfig:"REDIRECT_HOST" default:"localhost"`
	ServerPort uint32 `envconfig:"SERVER_PORT" default:"10000"`
}