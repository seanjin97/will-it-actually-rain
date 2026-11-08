package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"will-it-actually-rain/internal/geo"
)

func TestRun(t *testing.T) {
	var output bytes.Buffer
	search := func(postalCode string) (geo.LatLon, error) {
		if postalCode != "012345" {
			t.Fatalf("postal code = %q", postalCode)
		}
		return geo.LatLon{Latitude: 1.3003, Longitude: 103.8372}, nil
	}
	if err := run([]string{"012345"}, &output, search); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "S79\tSomerset Road\nS111\tScotts Road\nS222\tHenderson Road\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestRunErrorsAndNoStations(t *testing.T) {
	var output bytes.Buffer
	search := func(string) (geo.LatLon, error) { return geo.LatLon{}, errors.New("lookup failed") }
	if err := run(nil, &output, search); err == nil || !strings.Contains(err.Error(), "usage:") {
		t.Errorf("missing argument error = %v", err)
	}
	if err := run([]string{"123456"}, &output, search); err == nil || !strings.Contains(err.Error(), "lookup failed") {
		t.Errorf("lookup error = %v", err)
	}
	search = func(string) (geo.LatLon, error) { return geo.LatLon{}, nil }
	if err := run([]string{"123456"}, &output, search); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "No weather stations found within 3 km") {
		t.Errorf("output = %q", output.String())
	}
}
