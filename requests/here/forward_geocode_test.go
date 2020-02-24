package requests_test

import (
	"testing"

	"github.com/manfredkogler/dist-calc/models"
	requests "github.com/manfredkogler/dist-calc/requests/here"
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
	hereService := requests.HereGeoQuery{}

	routeInfo := hereService.CalculateRoute(Locations[Lambertgasse], Locations[Schottenring])
	routeInfo = hereService.CalculateRoute(Locations[Schottenring], Locations[Lambertgasse])
	_ = routeInfo

	location := hereService.ForwardGeocode("Schottenring 1 Wien")
	location = hereService.ForwardGeocode("Lambertgasse 4 Wien")
	_ = location
}
