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

// https://api-geocode.service.distancematrix.ai/maps/api/geocode/json?address=Lambertgasse+4,+Wien&language=de&key=YOUR_API_KEY
const (
	distancematrixaiAPIgeocodeURL       = "https://api-geocode.service.distancematrix.ai/maps/api/geocode/json?"
	distancematrixaiAPIgeocodeParamsURL = "address=%s&language=de"
)

var fowardGeocodeReqURL = distancematrixaiAPIgeocodeURL + distancematrixaiAPIgeocodeParamsURL + distanceMatrixAIendingCredentials

// ForwardGeocode returns the geocode for a given address specified as "searchString" (any string including whitespaces)
func (h DistancematrixaiGeoQuery) ForwardGeocode(searchString string) models.Loc {
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
	address := b.Result[0].FormattedAddress
	location := b.Result[0].Geometry.Location
	return models.Loc{
		Addr: address,
		Lat:  floatToString(location.Lat),
		Lng:  floatToString(location.Lng),
	}
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}
