package geo

// LatLon represents a position in WGS 84 (EPSG:4326), in decimal degrees.
// Latitude is positive north; longitude is positive east.
type LatLon struct {
	Latitude  float64
	Longitude float64
}
