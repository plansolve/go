// Package retry retries the polling calls of the solver clients on failures
// that may pass on a later attempt: a connection error or a 408, 409, 429 or
// 5xx answer, up to 3 times with 2, 4 and 8 second backoff, or after the delay
// a Retry-After header asks for (at most MaxRetryAfter). A request timeout is
// not retried.
package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// MaxRetries is the number of retries after the first attempt.
const MaxRetries = 3

// MaxRetryAfter caps the delay a Retry-After header can ask for.
const MaxRetryAfter = 60 * time.Second

// Backoff is the delay before retry number attempt (1-based) when the response
// carries no Retry-After header. A variable so tests can avoid real delays.
var Backoff = func(attempt int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt))) * time.Second
}

// StatusError is a non-2xx API response. Its message is the one the clients
// have always returned: "API error: status <code>: <message>".
type StatusError struct {
	StatusCode int
	Message    string
	// RetryAfter is the delay the Retry-After header asked for; see HasRetryAfter.
	RetryAfter time.Duration
	// HasRetryAfter reports whether the response carried a valid Retry-After header.
	HasRetryAfter bool
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("API error: status %d: %s", e.StatusCode, e.Message)
}

// NewStatusError builds the error for resp, whose body has been read as message.
func NewStatusError(resp *http.Response, message string) *StatusError {
	retryAfter, ok := parseRetryAfter(resp.Header.Get("Retry-After"), time.Now())
	return &StatusError{
		StatusCode:    resp.StatusCode,
		Message:       message,
		RetryAfter:    retryAfter,
		HasRetryAfter: ok,
	}
}

// parseRetryAfter reads a Retry-After value, a non-negative number of seconds
// or an HTTP date, capped at MaxRetryAfter, and whether it held a valid one. A
// date in the past reads as no delay.
func parseRetryAfter(value string, now time.Time) (time.Duration, bool) {
	if value == "" {
		return 0, false
	}
	if isDigits(value) {
		seconds, err := strconv.ParseUint(value, 10, 64)
		if err != nil || seconds > uint64(MaxRetryAfter/time.Second) {
			return MaxRetryAfter, true
		}
		return time.Duration(seconds) * time.Second, true
	}
	if at, err := http.ParseTime(value); err == nil {
		if !at.After(now) {
			return 0, true
		}
		if delay := at.Sub(now); delay < MaxRetryAfter {
			return delay, true
		}
		return MaxRetryAfter, true
	}
	return 0, false
}

func isDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// IsTransient reports whether err may pass on a later attempt: a 5xx, 408, 409
// or 429 response, or a connection error. A request timeout and a cancelled or
// expired ctx are not retried.
func IsTransient(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	var statusErr *StatusError
	if errors.As(err, &statusErr) {
		code := statusErr.StatusCode
		return code >= 500 || code == http.StatusRequestTimeout || code == http.StatusConflict || code == http.StatusTooManyRequests
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return !urlErr.Timeout()
	}
	return false
}

// Do calls fn and retries it up to MaxRetries times while it fails with a
// transient error, waiting for Retry-After when the response sent one and for
// Backoff otherwise.
func Do[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	for attempt := 1; ; attempt++ {
		result, err := fn()
		if err == nil || attempt > MaxRetries || !IsTransient(ctx, err) {
			return result, err
		}

		delay := Backoff(attempt)
		var statusErr *StatusError
		if errors.As(err, &statusErr) && statusErr.HasRetryAfter {
			delay = min(max(statusErr.RetryAfter, 0), MaxRetryAfter)
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			var zero T
			return zero, ctx.Err()
		case <-timer.C:
		}
	}
}
