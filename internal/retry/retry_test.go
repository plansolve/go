package retry

import (
	"context"
	"errors"
	"math"
	"net/url"
	"testing"
	"time"
)

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestIsTransient(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"503", &StatusError{StatusCode: 503}, true},
		{"500", &StatusError{StatusCode: 500}, true},
		{"429", &StatusError{StatusCode: 429}, true},
		{"408", &StatusError{StatusCode: 408}, true},
		{"409", &StatusError{StatusCode: 409}, true},
		{"400", &StatusError{StatusCode: 400}, false},
		{"401", &StatusError{StatusCode: 401}, false},
		{"404", &StatusError{StatusCode: 404}, false},
		{"422", &StatusError{StatusCode: 422}, false},
		{"connection reset", &url.Error{Op: "Get", URL: "x", Err: errors.New("connection reset")}, true},
		{"client timeout", &url.Error{Op: "Get", URL: "x", Err: timeoutError{}}, false},
		{"decode error", errors.New("failed to decode response"), false},
	}
	for _, c := range cases {
		if got := IsTransient(ctx, c.err); got != c.want {
			t.Errorf("%s: IsTransient = %v, want %v", c.name, got, c.want)
		}
	}

	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if IsTransient(cancelled, &StatusError{StatusCode: 503}) {
		t.Error("a cancelled context must not be retried")
	}
}

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, 11, 2, 9, 0, 0, 0, time.UTC)
	valid := map[string]time.Duration{
		"3":                             3 * time.Second,
		"0":                             0,
		"Mon, 02 Nov 2026 09:00:05 GMT": 5 * time.Second,
		"Mon, 02 Nov 2026 08:00:00 GMT": 0,
		"60":                            60 * time.Second,
		"3600":                          MaxRetryAfter,
		"99999999999999999999999":       MaxRetryAfter,
		"Mon, 02 Nov 2026 10:00:00 GMT": MaxRetryAfter,
	}
	for v, want := range valid {
		if got, ok := parseRetryAfter(v, now); !ok || got != want {
			t.Errorf("%q: got %v, %v, want %v, true", v, got, ok, want)
		}
	}
	for _, v := range []string{"", "-1", "1.5", "soon"} {
		if got, ok := parseRetryAfter(v, now); ok || got != 0 {
			t.Errorf("%q: got %v, %v, want 0, false", v, got, ok)
		}
	}
}

func TestDoRetriesTransientFailuresAndHonorsRetryAfter(t *testing.T) {
	var delays []time.Duration
	saved := Backoff
	Backoff = func(attempt int) time.Duration {
		delays = append(delays, saved(attempt))
		return 0
	}
	defer func() { Backoff = saved }()

	calls := 0
	got, err := Do(context.Background(), func() (string, error) {
		calls++
		if calls == 1 {
			return "", &StatusError{StatusCode: 503}
		}
		return "ok", nil
	})
	if err != nil || got != "ok" || calls != 2 {
		t.Fatalf("got %q, %v after %d calls", got, err, calls)
	}
	if len(delays) != 1 || delays[0] != 2*time.Second {
		t.Errorf("expected one 2 s backoff, got %v", delays)
	}

	calls = 0
	start := time.Now()
	_, err = Do(context.Background(), func() (string, error) {
		calls++
		if calls == 1 {
			return "", &StatusError{StatusCode: 429, RetryAfter: 50 * time.Millisecond, HasRetryAfter: true}
		}
		return "ok", nil
	})
	if err != nil || calls != 2 {
		t.Fatalf("429: %v after %d calls", err, calls)
	}
	if elapsed := time.Since(start); elapsed < 50*time.Millisecond {
		t.Errorf("Retry-After not honored: retried after %v", elapsed)
	}
}

func TestDoGivesUpAfterMaxRetriesAndOnPermanentErrors(t *testing.T) {
	saved := Backoff
	Backoff = func(int) time.Duration { return 0 }
	defer func() { Backoff = saved }()

	calls := 0
	_, err := Do(context.Background(), func() (int, error) {
		calls++
		return 0, &StatusError{StatusCode: 503}
	})
	if err == nil || calls != MaxRetries+1 {
		t.Errorf("expected %d calls and an error, got %d calls, err %v", MaxRetries+1, calls, err)
	}

	calls = 0
	_, err = Do(context.Background(), func() (int, error) {
		calls++
		return 0, &StatusError{StatusCode: 404}
	})
	if err == nil || calls != 1 {
		t.Errorf("404 must not be retried: %d calls, err %v", calls, err)
	}
}

func TestDoCapsRetryAfter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := Do(ctx, func() (int, error) {
		return 0, &StatusError{StatusCode: 429, RetryAfter: time.Duration(math.MaxInt64), HasRetryAfter: true}
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected the wait to last until the context expired, got %v after %v", err, time.Since(start))
	}
}
