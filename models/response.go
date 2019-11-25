package models

// CalculateRouteResponse embodies the interesting parts of the calculateroute request
type CalculateRouteResponse struct {
	Response cRresponse `json:"response"`
}

type cRresponse struct {
	Route []route `json:"route"`
}

type route struct {
	Summary summary `json:"summary"`
}

type summary struct {
	Distance   int `json:"distance"`
	TravelTime int `json:"travelTime"`
}

// ForwardGeocoderResponse embodies the interesting parts of the forward geocode request
type ForwardGeocoderResponse struct {
	Response fGCresponse `json:"Response"`
}

type fGCresponse struct {
	View []view `json:"View"`
}

type view struct {
	Result []result `json:"Result"`
}

type result struct {
	Location location `json:"Location"`
}

type location struct {
	NavigationPosition []navigationPosition `json:"NavigationPosition"`
	Address            address              `json:"Address"`
}

type navigationPosition struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type address struct {
	Label string `json:"Label"`
}
