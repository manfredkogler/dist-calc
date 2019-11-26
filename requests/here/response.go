package requests

// calculateRouteResponse embodies the interesting parts of the calculateroute request
type calculateRouteResponse struct {
	Response struct {
		Route []struct {
			Summary struct {
				Distance   int `json:"distance"`
				TravelTime int `json:"travelTime"`
			} `json:"summary"`
		} `json:"route"`
	} `json:"response"`
}

// forwardGeocodeResponse embodies the interesting parts of the forward geocode request
type forwardGeocodeResponse struct {
	Response struct {
		View []struct {
			Result []struct {
				Location struct {
					NavigationPosition []struct {
						Latitude  float64 `json:"Latitude"`
						Longitude float64 `json:"Longitude"`
					} `json:"NavigationPosition"`
					Address struct {
						Label string `json:"Label"`
					} `json:"Address"`
				} `json:"Location"`
			} `json:"Result"`
		} `json:"View"`
	} `json:"Response"`
}
