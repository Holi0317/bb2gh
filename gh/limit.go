package gh

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/sirupsen/logrus"
)

const retry = 3

func getRetryAfter(resp *http.Response) time.Duration {
	after := resp.Header.Get("Retry-After")
	if after == "" {
		return 30 * time.Second
	}

	parsed, err := strconv.Atoi(after)
	if err != nil {
		return 30 * time.Second
	}

	return time.Duration(parsed) * time.Second
}

func isRateLimit(err error) (time.Duration, bool) {
	if err == nil {
		return 0, false
	}

	if rlerr, ok := err.(*github.RateLimitError); ok {
		after := getRetryAfter(rlerr.Response)
		logrus.WithField("after", after).WithError(rlerr).Warn("Got hard rate limit error")
		return after, true
	}

	if aberr, ok := err.(*github.AbuseRateLimitError); ok {
		after := getRetryAfter(aberr.Response)
		logrus.WithField("after", after).WithError(aberr).Warn("Got hard rate limit error")
		return after, true
	}

	return 0, false
}
