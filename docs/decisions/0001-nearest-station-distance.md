# 1. Nearest station distance

- **Status:** Accepted
- **Context:** Find nearby stations for a WGS 84 (EPSG:4326) location.
- **Decision:** Compare WGS 84 ellipsoidal geodesic distances using GeographicLib's Karney algorithm; return up to three IDs within 3 km, closest first.
- **Rationale:** The ellipsoid matches the coordinate system more closely than a spherical haversine calculation. A linear scan is suitable for the small station list.
- **Radius:** 3 km is a provisional local cutoff, informed by a [Singapore rainfall study](https://journals.ametsoc.org/view/journals/hydr/17/5/jhm-d-15-0100_1.xml) using that radius to fill daily gauge data; validate it against the intended rainfall estimate.
- **Consequences:** Coordinate accuracy and rainfall variability still limit the result. No nearby stations returns an empty list; ties retain input order.
