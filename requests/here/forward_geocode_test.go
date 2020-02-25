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
			Addr: "Lambertgasse 4, 1160 Wien, Österreich",
			Lat:  "48.211820",
			Lng:  "16.319660"},
	},
	{
		"Schottenring 1 Wien",
		models.Loc{
			Addr: "Schottenring 1, 1010 Wien, Österreich",
			Lat:  "48.214180",
			Lng:  "16.362830"},
	},
}

// test data END

func TestForwardGeocode(t *testing.T) {
	geoService := HereGeoQuery{}

	for _, td := range testDatas {
		checkForwardGeocode(t, geoService, td)
	}
}

func checkForwardGeocode(t *testing.T, s HereGeoQuery, td testData) {
	got := s.ForwardGeocode(td.addr)
	if got != td.loc {
		t.Errorf("ForwardGeocode() = %v; want %v", got, td.loc)
	}
}

func TestCalculateRoute(t *testing.T) {
	geoService := HereGeoQuery{}

	want := models.RouteInfo{
		Distance:   3736,
		TravelTime: 583,
	}

	got := geoService.CalculateRoute(testDatas[0].loc, testDatas[1].loc)
	if got != want {
		t.Errorf("CalculateRoute() = %v; want %v", got, want)
	}
}
