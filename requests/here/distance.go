package requests

import (
	"bytes"
	"dist-calc/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// curl 'https://route.api.here.com/routing/7.2/calculateroute.json?app_id=lD6sZ3QeKG552sEIkRVn&app_code=TbdJaZQdA7QxIIc3Qj--7A&waypoint0=geo!48.211836,16.319760&waypoint1=geo!48.215239,16.365305&representation=overview&mode=fastest;car;traffic:enabled'
const (
	hereAPIroutingURL           = "https://route.api.here.com/routing/7.2/calculateroute.json?"
	hereAPIroutingTailParamsURL = "waypoint0=geo!%s,%s&waypoint1=geo!%s,%s&representation=overview&mode=fastest;car;traffic:disabled"
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

// CachedCalculateRouteClosure returns a cached version of CalculateRoute
func (h HereGeoQuery) CachedCalculateRouteClosure() (f func(models.Loc, models.Loc) (models.RouteInfo, bool)) {
	// The cache
	var routeInfoMap = map[string]models.RouteInfo{}

	f = func(from models.Loc, to models.Loc) (models.RouteInfo, bool) {
		route := from.Addr + " -> " + to.Addr
		routeInfo, ok := routeInfoMap[route]
		if ok {
			return routeInfo, ok
		}
		routeInfo = h.CalculateRoute(from, to)
		routeInfoMap[route] = routeInfo
		return routeInfo, ok
	}
	return
}
