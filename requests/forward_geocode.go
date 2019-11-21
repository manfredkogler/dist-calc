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

// curl 'https://geocoder.api.here.com/6.2/geocode.json?app_id=lD6sZ3QeKG552sEIkRVn&app_code=TbdJaZQdA7QxIIc3Qj--7A&searchtext=Schottenring+1+Wien'
const (
	hereAPIgeocoderURL           = "https://geocoder.api.here.com/6.2/geocode.json?"
	hereAPIgeocoderTailParamsURL = "searchtext=%s"
)

var fowardGeocoderReqURL = hereAPIgeocoderURL + hereAPIstartingCredentials + hereAPIgeocoderTailParamsURL

// test data BEGIN

// test data END

// ForwardGeocode returns the geocode for a given address specified as "searchString" (strings separated by blanks)
func ForwardGeocode(searchString string) models.Loc {
	reqString := fmt.Sprintf(fowardGeocoderReqURL, url.QueryEscape(searchString))
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

	var b models.ForwardGeocoderResponse = models.ForwardGeocoderResponse{}
	jd := json.NewDecoder(response.Body)
	jd.Decode(&b)
	resultLocation := b.Response.View[0].Result[0].Location
	return models.Loc{
		Addr: resultLocation.Address.Label,
		Lat:  floatToString(resultLocation.NavigationPosition[0].Latitude),
		Lng:  floatToString(resultLocation.NavigationPosition[0].Longitude),
	}
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}
