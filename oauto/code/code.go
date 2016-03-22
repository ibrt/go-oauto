package main

import (
	"regexp"
	"time"
)

type Code interface {
	Receive(pattern *regexp.Regexp, after time.Time, retryCount int, retryDelay time.Duration) (string, error)
}
