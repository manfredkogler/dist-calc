package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/manfredkogler/dist-calc/models"
)

// https://geocoder.ls.hereapi.com/6.2/geocode.json?apiKey={YOUR_API_KEY}&searchtext=Schottenring+1+Wien&language=de-de
const (
	hereAPIgeocodeURL           = "https://geocoder.ls.hereapi.com/6.2/geocode.json?"
	hereAPIgeocodeTailParamsURL = "searchtext=%s&language=de-de"
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
	// The dot as decimal point must not be changed as the maps services need it like that.
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}
