package gh

import (
	"time"

	"github.com/google/go-github/v48/github"
)

const maxProcess = 10
const retry = 3

func sleepTime(i int) time.Duration {
	return time.Duration(30*(i+1)) * time.Second
}

func isRateLimit(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(*github.RateLimitError); ok {
		return true
	}

	if _, ok := err.(*github.AbuseRateLimitError); ok {
		return true
	}

	return false
}
