package onemap

import (
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"will-it-actually-rain/internal/geo"
)

// Search returns the coordinates of the first OneMap result for a postal code.
func Search(postalCode string) (geo.LatLon, error) {
	if len(postalCode) != 6 {
		return geo.LatLon{}, fmt.Errorf("invalid postal code: %q (must be six digits)", postalCode)
	}
	for _, digit := range postalCode {
		if digit < '0' || digit > '9' {
			return geo.LatLon{}, fmt.Errorf("invalid postal code: %q (must be six digits)", postalCode)
		}
	}
	token, err := AccessToken()
	if err != nil {
		return geo.LatLon{}, fmt.Errorf("search OneMap: %w", err)
	}
	query := url.Values{"searchVal": {postalCode}, "returnGeom": {"Y"}, "getAddrDetails": {"N"}}
	req, err := http.NewRequest(http.MethodGet, "https://www.onemap.gov.sg/api/common/elastic/search?"+query.Encode(), nil)
	if err != nil {
		return geo.LatLon{}, fmt.Errorf("create OneMap search request: %w", err)
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return geo.LatLon{}, fmt.Errorf("search OneMap: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return geo.LatLon{}, fmt.Errorf("search OneMap: HTTP %s", resp.Status)
	}
	var result struct {
		Error   string `json:"error"`
		Results []struct {
			Latitude  string `json:"LATITUDE"`
			Longitude string `json:"LONGITUDE"`
		} `json:"results"`
	}
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return geo.LatLon{}, fmt.Errorf("decode OneMap search response: %w", err)
	}
	if result.Error != "" {
		return geo.LatLon{}, fmt.Errorf("search OneMap: %s", result.Error)
	}
	if len(result.Results) == 0 {
		return geo.LatLon{}, fmt.Errorf("no OneMap results for postal code %s", postalCode)
	}
	lat, err := strconv.ParseFloat(result.Results[0].Latitude, 64)
	if err != nil {
		return geo.LatLon{}, fmt.Errorf("parse OneMap latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(result.Results[0].Longitude, 64)
	if err != nil {
		return geo.LatLon{}, fmt.Errorf("parse OneMap longitude: %w", err)
	}
	return geo.LatLon{Latitude: lat, Longitude: lon}, nil
}
