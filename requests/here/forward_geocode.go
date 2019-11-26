package requests

import (
	"bytes"
	"dist-calc/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// curl "https://geocoder.api.here.com/6.2/geocode.json?app_id={YOUR_APP_ID}&app_code={YOUR_APP_CODE}&searchtext=Schottenring+1+Wien"
const (
	hereAPIgeocodeURL           = "https://geocoder.api.here.com/6.2/geocode.json?"
	hereAPIgeocodeTailParamsURL = "searchtext=%s"
)

var fowardGeocodeReqURL = hereAPIgeocodeURL + hereAPIstartingCredentials + hereAPIgeocodeTailParamsURL

// ForwardGeocode returns the geocode for a given address specified as "searchString" (any string including whitespaces)
func (h HereGeoQuery) ForwardGeocode(searchString string) models.Loc {
	reqString := fmt.Sprintf(fowardGeocodeReqURL, url.QueryEscape(searchString))
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

	var b forwardGeocodeResponse
	jd := json.NewDecoder(response.Body)
	jd.Decode(&b)
	location := b.Response.View[0].Result[0].Location
	return models.Loc{
		Addr: location.Address.Label,
		Lat:  floatToString(location.NavigationPosition[0].Latitude),
		Lng:  floatToString(location.NavigationPosition[0].Longitude),
	}
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}