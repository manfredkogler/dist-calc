package models

// CalculateRouteResponse embodies the interesting parts of the calculateroute request
type CalculateRouteResponse struct {
	Response struct {
		Route []struct {
			Summary struct {
				Distance   int `json:"distance"`
				TravelTime int `json:"travelTime"`
			} `json:"summary"`
		} `json:"route"`
	} `json:"response"`
}

// ForwardGeocoderResponse embodies the interesting parts of the forward geocode request
type ForwardGeocoderResponse struct {
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
