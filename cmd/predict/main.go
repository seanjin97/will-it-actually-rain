package main

import (
	"fmt"
	"io"
	"os"

	"will-it-actually-rain/internal/geo"
	"will-it-actually-rain/internal/integrations/nea"
	"will-it-actually-rain/internal/integrations/onemap"
	"will-it-actually-rain/internal/prediction"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, onemap.Search); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer, search func(string) (geo.LatLon, error)) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: predict <six-digit postal code>")
	}
	location, err := search(args[0])
	if err != nil {
		return fmt.Errorf("find postal code: %w", err)
	}
	ids := prediction.GetNearestWeatherStations(location, nea.Stations)
	if len(ids) == 0 {
		_, err = fmt.Fprintf(out, "No weather stations found within %d km of postal code %s.\n", prediction.RainfallStationRadiusMeters/1000, args[0])
		return err
	}
	stations := make(map[string]nea.Station, len(nea.Stations))
	for _, station := range nea.Stations {
		stations[station.ID] = station
	}
	for _, id := range ids {
		if _, err := fmt.Fprintf(out, "%s\t%s\n", id, stations[id].Name); err != nil {
			return err
		}
	}
	return nil
}
