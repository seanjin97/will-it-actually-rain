package prediction

import (
	"sort"

	"github.com/pymaxion/geographiclib-go/v2/geodesic"

	"will-it-actually-rain/internal/geo"
	"will-it-actually-rain/internal/integrations/nea"
)

// RainfallStationRadiusMeters is a provisional limit for local rainfall estimates.
const RainfallStationRadiusMeters = 3000

// GetNearestWeatherStations returns up to three station IDs within the radius,
// ordered by distance. Equidistant stations retain their input order.
func GetNearestWeatherStations(latlon geo.LatLon, sensors []nea.Station) []string {
	type candidate struct {
		id       string
		distance float64
	}
	candidates := make([]candidate, 0, len(sensors))
	for _, sensor := range sensors {
		distance := geodesic.WGS84.Inverse(latlon.Latitude, latlon.Longitude, sensor.LatLon.Latitude, sensor.LatLon.Longitude).S12
		if distance <= RainfallStationRadiusMeters {
			candidates = append(candidates, candidate{id: sensor.ID, distance: distance})
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].distance < candidates[j].distance })
	if len(candidates) > 3 {
		candidates = candidates[:3]
	}
	ids := make([]string, 0, len(candidates))
	for _, station := range candidates {
		ids = append(ids, station.id)
	}
	return ids
}
