package prediction

import (
	"reflect"
	"testing"

	"github.com/pymaxion/geographiclib-go/v2/geodesic"

	"will-it-actually-rain/internal/geo"
	"will-it-actually-rain/internal/integrations/nea"
)

func TestGetNearestWeatherStations(t *testing.T) {
	position := func(meters float64) geo.LatLon {
		point := geodesic.WGS84.Direct(0, 0, 90, meters)
		return geo.LatLon{Latitude: point.Lat2, Longitude: point.Lon2}
	}
	tests := []struct {
		name     string
		location geo.LatLon
		stations []nea.Station
		want     []string
	}{
		{name: "empty", want: []string{}},
		{name: "outside radius", stations: []nea.Station{{ID: "far", LatLon: position(3001)}}, want: []string{}},
		{name: "sort limit and radius", stations: []nea.Station{
			{ID: "fourth", LatLon: position(2999)},
			{ID: "outside", LatLon: position(3001)},
			{ID: "second", LatLon: position(1000)},
			{ID: "first", LatLon: position(100)},
			{ID: "third", LatLon: position(2000)},
		}, want: []string{"first", "second", "third"}},
		{name: "just inside radius", stations: []nea.Station{{ID: "near", LatLon: position(2999)}}, want: []string{"near"}},
		{name: "ellipsoidal distance", stations: []nea.Station{
			{ID: "east", LatLon: geo.LatLon{Longitude: 0.00996}},
			{ID: "north", LatLon: geo.LatLon{Latitude: 0.01}},
		}, want: []string{"north", "east"}},
		{name: "across date line", location: geo.LatLon{Longitude: 179.999}, stations: []nea.Station{
			{ID: "near", LatLon: geo.LatLon{Longitude: -179.999}},
		}, want: []string{"near"}},
		{name: "tie retains input order", stations: []nea.Station{
			{ID: "first", LatLon: geo.LatLon{Latitude: 0.01}},
			{ID: "second", LatLon: geo.LatLon{Latitude: -0.01}},
		}, want: []string{"first", "second"}},
		{name: "cached station", location: geo.LatLon{Latitude: 1.3003, Longitude: 103.8372}, stations: nea.Stations, want: []string{"S79", "S111", "S222"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNearestWeatherStations(tt.location, tt.stations); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNearestWeatherStations() = %v, want %v", got, tt.want)
			}
		})
	}
}
