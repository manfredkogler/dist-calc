package requests

// distancematrixResponse embodies the interesting parts of the distancematrix request
type distancematrixResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Value int `json:"value"`
			} `json:"distance"`
			Duration struct {
				Value int `json:"value"`
			} `json:"duration"`
		} `json:"elements"`
	} `json:"rows"`
}

// forwardGeocodeResponse embodies the interesting parts of the forward geocode request
type forwardGeocodeResponse struct {
	Result []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"result"`
}
