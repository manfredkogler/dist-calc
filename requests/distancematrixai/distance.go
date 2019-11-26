package requests

import (
	"bytes"
	"dist-calc/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// https://api.distancematrix.ai/maps/api/distancematrix/json?origins=48.211836,16.319760&destinations=48.215239,16.365305&key=<your_access_token>
const (
	distancematrixaiAPIdistancematrixURL       = "https://api.distancematrix.ai/maps/api/distancematrix/json?"
	distancematrixaiAPIdistancematrixParamsURL = "origins=%s,%s&destinations=%s,%s"
)

var distanceReqURL = distancematrixaiAPIdistancematrixURL + distancematrixaiAPIdistancematrixParamsURL + distanceMatrixAIendingCredentials

// CalculateRoute calculates and returns the route info from "from" to "to"
func (h DistancematrixaiGeoQuery) CalculateRoute(from models.Loc, to models.Loc) models.RouteInfo {
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

	var b distancematrixResponse
	jd := json.NewDecoder(response.Body)
	jd.Decode(&b)
	distance := b.Rows[0].Elements[0].Distance.Value
	duration := b.Rows[0].Elements[0].Duration.Value
	return models.RouteInfo{Distance: distance, TravelTime: duration}
}
