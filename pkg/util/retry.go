package util

import (
	"strings"
	"time"

	"github.com/avast/retry-go"
)

// add some retry logic here as some registries can be flaky
func RetryOn500(retryable func() error) error {
	return RetryOnFunc(retryable, func(err error) bool {
		return strings.Contains(err.Error(), "500 Internal Server Error")
	},
		retry.Attempts(4),
		retry.Delay(250*time.Millisecond),
	)
}

// add some retry logic here as some registries can be flaky
func RetryOnFunc(retryable func() error, retryIf func(err error) bool, opts ...retry.Option) error {
	return retry.Do(retryable,
		append(opts, retry.RetryIf(retryIf))...,
	)
}
