package onemap

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"will-it-actually-rain/internal/geo"
)

func TestSearch(t *testing.T) {
	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
		tokenCache.accessToken, tokenCache.expiresAt = "", 0
	}()
	tokenCache.accessToken, tokenCache.expiresAt = "test-token", time.Now().Add(time.Hour).Unix()
	body := `{"results":[{"LATITUDE":"1.30743547948389","LONGITUDE":"103.854713903431"}]}`
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/common/elastic/search" || r.URL.Query().Get("searchVal") != "012345" || r.URL.Query().Get("returnGeom") != "Y" || r.URL.Query().Get("getAddrDetails") != "N" || r.Header.Get("Authorization") != "test-token" {
			t.Errorf("unexpected request: %s %s, authorization %q", r.Method, r.URL, r.Header.Get("Authorization"))
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})
	got, err := Search("012345")
	want := (geo.LatLon{Latitude: 1.30743547948389, Longitude: 103.854713903431})
	if err != nil || got != want {
		t.Fatalf("Search() = %+v, %v", got, err)
	}
	for _, response := range []string{`{"results":[]}`, `{"error":"invalid token"}`, `{"results":[{"LATITUDE":"bad","LONGITUDE":"103"}]}`, `{`} {
		body = response
		if _, err := Search("012345"); err == nil {
			t.Errorf("Search() expected error for %s", response)
		}
	}
}

func TestSearchRejectsInvalidPostalCodes(t *testing.T) {
	for _, postalCode := range []string{"", "12345", "1234567", "12a456", "-12345", "１２３４５６"} {
		if _, err := Search(postalCode); err == nil {
			t.Errorf("Search(%q) expected an error", postalCode)
		}
	}
}
