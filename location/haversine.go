// This package has been copied from https://github.com/umahmood/haversine/blob/master/haversine.go
// to avoid an import. The funcs and types have been lowercased to avoid
// making them public.
//
// umahmood's haversine package is released under the MIT licence.
package location

import (
	"math"
)

const (
	earthRadiusMi = 3958 // radius of the earth in miles.
	earthRaidusKm = 6371 // radius of the earth in kilometers.
)

// coord represents a geographic coordinate.
type coord struct {
	Lat float64
	Lon float64
}

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// haversineDistance calculates the shortest path between two
// coordinates on the surface of the Earth. This function returns two
// units of measure, the first is the distance in miles, the second is
// the distance in kilometers.
func haversineDistance(p, q coord) (mi, km float64) {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	mi = c * earthRadiusMi
	km = c * earthRaidusKm

	return mi, km
}
