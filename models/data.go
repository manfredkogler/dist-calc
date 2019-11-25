package models

// Loc stores location data
type Loc struct {
	Addr, Lat, Lng string
}

// RouteInfo stores route information
type RouteInfo struct {
	Distance   int
	TravelTime int
}
