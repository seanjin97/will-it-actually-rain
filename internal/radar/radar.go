package radar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const Endpoint = "https://api-open.data.gov.sg/v2/real-time/api/weather-radar-images/70km"

type radarResponse struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"errorMsg"`
	Data     struct {
		Records []struct {
			Timestamp time.Time `json:"timestamp"`
			Image     struct {
				URL string `json:"url"`
			} `json:"image"`
		} `json:"records"`
	} `json:"data"`
}

// CacheImage retrieves the image for a radar slot unless it is already cached.
func CacheImage(client *http.Client, endpoint, dir string, slot time.Time) error {
	filename := filepath.Join(dir, "rain_70km_"+slot.Format("20060102150405")+".png")
	if info, err := os.Stat(filename); err == nil && info.Mode().IsRegular() && info.Size() > 0 {
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("check cache: %w", err)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse radar endpoint: %w", err)
	}
	query := u.Query()
	query.Set("date", slot.Format("2006-01-02T15:04:05"))
	u.RawQuery = query.Encode()
	res, err := client.Get(u.String())
	if err != nil {
		return fmt.Errorf("request radar metadata: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("radar metadata returned HTTP %d", res.StatusCode)
	}
	var data radarResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return fmt.Errorf("decode radar metadata: %w", err)
	}
	if data.Code != 0 {
		return fmt.Errorf("radar API error %d: %s", data.Code, data.ErrorMsg)
	}
	for _, record := range data.Data.Records {
		if !record.Timestamp.Equal(slot) {
			continue
		}
		if record.Image.URL == "" {
			return fmt.Errorf("radar record has no image URL")
		}
		return downloadRadarImage(client, record.Image.URL, dir, filename)
	}
	return fmt.Errorf("no radar image for %s", slot.Format(time.RFC3339))
}

func downloadRadarImage(client *http.Client, imageURL, dir, filename string) error {
	res, err := client.Get(imageURL)
	if err != nil {
		return fmt.Errorf("download radar image: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("radar image returned HTTP %d", res.StatusCode)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create radar cache: %w", err)
	}
	file, err := os.CreateTemp(dir, ".radar-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary radar image: %w", err)
	}
	defer os.Remove(file.Name())
	written, copyErr := io.Copy(file, res.Body)
	closeErr := file.Close()
	if copyErr != nil {
		return fmt.Errorf("save radar image: %w", copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close radar image: %w", closeErr)
	}
	if written == 0 {
		return fmt.Errorf("radar image is empty")
	}
	if err := os.Rename(file.Name(), filename); err != nil {
		return fmt.Errorf("cache radar image: %w", err)
	}
	return nil
}
