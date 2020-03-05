package util

import (
	"strings"
	"time"

	"github.com/avast/retry-go"
)

// add some retry logic here as some registries can be flaky
func RetryOn500(retryable func() error) error {
	return retry.Do(retryable,
		retry.Attempts(4),
		retry.Delay(250*time.Millisecond),
		retry.RetryIf(func(err error) bool {
			return strings.Contains(err.Error(), "500 Internal Server Error")
		}),
	)
}
