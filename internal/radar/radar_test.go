package radar

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}
}

func TestCacheRadarImage(t *testing.T) {
	slot := time.Date(2026, 11, 8, 23, 50, 0, 0, time.FixedZone("SGT", 8*60*60))
	requests := 0
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		if r.URL.Path == "/metadata" {
			if got := r.URL.Query().Get("date"); got != "2026-11-08T23:50:00" {
				t.Errorf("date = %q", got)
			}
			return response(200, `{"code":0,"data":{"records":[{"timestamp":"2026-11-08T23:50:00+08:00","image":{"url":"https://example.com/image?signature=test"}}]}}`), nil
		}
		return response(200, "image data"), nil
	})}
	dir := filepath.Join(t.TempDir(), "radars")
	for i := 0; i < 2; i++ {
		if err := CacheImage(client, "https://example.com/metadata", dir, slot); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(filepath.Join(dir, "rain_70km_20261108235000.png"))
	if err != nil || string(got) != "image data" {
		t.Fatalf("cached image = %q, error = %v", got, err)
	}
	if requests != 2 {
		t.Fatalf("got %d requests, want 2", requests)
	}
}

func TestCacheRadarImageRejectsFailedDownload(t *testing.T) {
	slot := time.Date(2026, 11, 8, 23, 50, 0, 0, time.FixedZone("SGT", 8*60*60))
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/image" {
			return response(403, "expired"), nil
		}
		return response(200, fmt.Sprintf(`{"code":0,"data":{"records":[{"timestamp":"2026-11-08T23:50:00+08:00","image":{"url":%q}}]}}`, "https://example.com/image")), nil
	})}
	dir := t.TempDir()
	if err := CacheImage(client, "https://example.com/metadata", dir, slot); err == nil {
		t.Fatal("expected download error")
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("cache entries = %v, error = %v", entries, err)
	}
}
