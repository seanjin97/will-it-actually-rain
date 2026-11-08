package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"will-it-actually-rain/internal/radar"
)

func main() {
	location, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		fmt.Fprintln(os.Stderr, "load Singapore time zone:", err)
		os.Exit(1)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	latest := time.Now().In(location).Truncate(5 * time.Minute)
	failed := false
	for i := 0; i < 12; i++ {
		slot := latest.Add(-time.Duration(i) * 5 * time.Minute)
		if err := radar.CacheImage(client, radar.Endpoint, "radars", slot); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", slot.Format(time.RFC3339), err)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
