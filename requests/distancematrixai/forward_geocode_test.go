package requests

import (
	"testing"

	"github.com/manfredkogler/dist-calc/models"
)

// test data BEGIN

// Test data
type testData struct {
	addr string
	loc  models.Loc
}

var testDatas []testData = []testData{
	{
		"Lambertgasse 4 Wien",
		models.Loc{
			Addr: "Lambertgasse 4 Wien",
			Lat:  "48.211824",
			Lng:  "16.319796"},
	},
	{
		"Schottenring 1 Wien",
		models.Loc{
			Addr: "Schottenring 1 Wien",
			Lat:  "48.215016",
			Lng:  "16.364444"},
	},
}

// test data END

func TestForwardGeocode(t *testing.T) {
	geoService := DistancematrixaiGeoQuery{}

	for _, td := range testDatas {
		checkForwardGeocode(t, geoService, td)
	}
}

func checkForwardGeocode(t *testing.T, s DistancematrixaiGeoQuery, td testData) {
	got := s.ForwardGeocode(td.addr)
	if got != td.loc {
		t.Errorf("ForwardGeocode() = %v; want %v", got, td.loc)
	}
}

func TestCalculateRoute(t *testing.T) {
	geoService := DistancematrixaiGeoQuery{}

	want := models.RouteInfo{
		Distance:   3688,
		TravelTime: 803,
	}

	got := geoService.CalculateRoute(testDatas[0].loc, testDatas[1].loc)
	if got != want {
		t.Errorf("CalculateRoute() = %v; want %v", got, want)
	}
}
