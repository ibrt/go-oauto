package main

import (
	"regexp"
	"time"
)

type Twilio struct {
	sid    string
	secret string
}

func NewTwilio(sid, secret string) *Twilio {
	return &Twilio{
		sid:    sid,
		secret: secret,
	}
}

func (t *Twilio) Receive(pattern *regexp.Regexp, after time.Time, retryCount int, retryDelay time.Duration) (string, error) {

	return "", nil
}
