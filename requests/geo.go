package requests

import "dist-calc/models"

// CachedCalculateRoute is cached version of CalculateRoute
// type CachedCalculateRoute func(models.Loc, models.Loc) (models.RouteInfo, bool)

// CachedForwardGeocode is cached version of ForwardGeocode
// type CachedForwardGeocode func(string) (models.Loc, bool)

// GeoQuery defines the interface for some geo queries
type GeoQuery interface {
	// CalculateRoute calculates and returns the route info from "from" to "to"
	CalculateRoute(from models.Loc, to models.Loc) models.RouteInfo
	// ForwardGeocode returns the geocode for a given address specified as "searchString" (any string including whitespaces)
	ForwardGeocode(searchString string) models.Loc

	// CachedCalculateRouteClosure returns a cached version of CalculateRoute
	CachedCalculateRouteClosure() func(models.Loc, models.Loc) (models.RouteInfo, bool)
	// CachedForwardGeocodeClosure returns a cached version of ForwardGeocode
	CachedForwardGeocodeClosure() func(string) (models.Loc, bool)
}
