package config

type Config struct {
	SeleniumURL  string `envconfig:"SELENIUM_URL" default:"http://localhost:4444/wd/hub"`
	RedirectHost string `envconfig:"REDIRECT_HOST" default:"localhost"`
	ServerPort   uint32 `envconfig:"SERVER_PORT" default:"10000"`
	TwilioSID    string `envconfig:"TWILIO_SID" default:""`
	TwilioSecret string `envconfig:"TWILIO_SECRET" default:""`
}
