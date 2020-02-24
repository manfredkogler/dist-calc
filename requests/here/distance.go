package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manfredkogler/dist-calc/models"
)

// https://route.ls.hereapi.com/routing/7.2/calculateroute.json?apiKey={YOUR_API_KEY}&waypoint0=geo!48.211836,16.319760&waypoint1=geo!48.215239,16.365305&representation=overview&mode=fastest;car;traffic:disabled&language=de-de
const (
	hereAPIroutingURL           = "https://route.ls.hereapi.com/routing/7.2/calculateroute.json?"
	hereAPIroutingTailParamsURL = "waypoint0=geo!%s,%s&waypoint1=geo!%s,%s&representation=overview&mode=fastest;car;traffic:disabled&language=de-de"
)

var distanceReqURL = hereAPIroutingURL + hereAPIstartingCredentials + hereAPIroutingTailParamsURL

// CalculateRoute calculates and returns the route info from "from" to "to"
func (h HereGeoQuery) CalculateRoute(from models.Loc, to models.Loc) models.RouteInfo {
	reqString := fmt.Sprintf(distanceReqURL, from.Lat, from.Lng, to.Lat, to.Lng)
	fmt.Println(reqString)
	response, err := http.Get(reqString)
	var data []byte
	if err != nil {
		fmt.Printf("The HTTP request failed with rror %s\n", err)
	} else {
		defer response.Body.Close()
		data, _ = ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

	// Restore the io.ReadCloser to its original state
	response.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	var b calculateRouteResponse
	jd := json.NewDecoder(response.Body)
	jd.Decode(&b)
	rs := b.Response.Route[0].Summary
	return models.RouteInfo{Distance: rs.Distance, TravelTime: rs.TravelTime}
}
