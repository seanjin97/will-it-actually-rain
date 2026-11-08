package onemap

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestAccessTokenCaching(t *testing.T) {
	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
		tokenCache.accessToken, tokenCache.expiresAt = "", 0
	}()
	tokenCache.accessToken, tokenCache.expiresAt = "", 0
	t.Setenv("ONE_MAP_EMAIL", "test@example.com")
	t.Setenv("ONE_MAP_PASSWORD", "password")

	calls := 0
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodPost || r.URL.String() != endpoint {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		}
		body := `{"access_token":"token","expiry_timestamp":"4102444800"}`
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})

	for i := 0; i < 2; i++ {
		if token, err := AccessToken(); err != nil || token != "token" {
			t.Fatalf("AccessToken() = %q, %v", token, err)
		}
	}
	if calls != 1 {
		t.Fatalf("got %d authentication calls, want 1", calls)
	}
	tokenCache.expiresAt = time.Now().Unix() - 1
	if _, err := AccessToken(); err != nil || calls != 2 {
		t.Fatalf("expired token: calls = %d, error = %v", calls, err)
	}
}
