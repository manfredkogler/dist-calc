package requests_test

import (
	"dist-calc/models"
	"dist-calc/requests"
	"testing"
)

// test data BEGIN

// Locations test data
var Locations []models.Loc = []models.Loc{
	{Addr: "Lambertgasse 4", Lat: "48.211836", Lng: "16.319760"},
	{Addr: "Schottenring 1", Lat: "48.215239", Lng: "16.365305"},
}

const (
	// Lambertgasse test data
	Lambertgasse = 0
	// Schottenring test data
	Schottenring = 1
)

// test data END

func TestForwardGeocode(t *testing.T) {
	routeInfo := requests.CalculateRoute(requests.Locations[requests.Lambertgasse], requests.Locations[requests.Schottenring])
	routeInfo = requests.CalculateRoute(requests.Locations[requests.Schottenring], requests.Locations[requests.Lambertgasse])
	_ = routeInfo

	location := requests.ForwardGeocode("Schottenring 1 Wien")
	location = requests.ForwardGeocode("Lambertgasse 4 Wien")
	_ = location
}
