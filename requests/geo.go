package requests

import "dist-calc/models"

// CachedCalculateRoute is cached version of CalculateRoute
// type CachedCalculateRoute func(models.Loc, models.Loc) (models.RouteInfo, bool)

// CachedForwardGeocode is cached version of ForwardGeocode
// type CachedForwardGeocode func(string) (models.Loc, bool)

// GeoQuery defines the interface for some geo queries
type GeoQuery interface {
	// ForwardGeocode returns the geocode for a given address specified as "searchString" (any string including whitespaces)
	ForwardGeocode(searchString string) models.Loc
	// CalculateRoute calculates and returns the route info from "from" to "to"
	CalculateRoute(from models.Loc, to models.Loc) models.RouteInfo
}

// CachedGeoQuery defines the interface for a cached version of GeoQuery
type CachedGeoQuery interface {
	// ForwardGeocode is a cached version of GeoQuery.ForwardGeocode. The bool returned signals if the returned entry is taken from the cache.
	ForwardGeocode(searchString string) (models.Loc, bool)
	// CalculateRoute is a cached version of GeoQuery.CalculateRoute. The bool returned signals if the returned entry is taken from the cache.
	CalculateRoute(from models.Loc, to models.Loc) (models.RouteInfo, bool)
}

// cachedGeoQueryImpl implements a cached version of GeoQuery
type cachedGeoQueryImpl struct {
	geoQuery             GeoQuery
	cachedForwardGeocode func(string) (models.Loc, bool)
	cachedCalculateRoute func(models.Loc, models.Loc) (models.RouteInfo, bool)
}

// CachedForwardGeocode ... see interface
func (c cachedGeoQueryImpl) ForwardGeocode(searchString string) (models.Loc, bool) {
	return c.cachedForwardGeocode(searchString)
}

// CachedCalculateRoute ... see interface
func (c cachedGeoQueryImpl) CalculateRoute(from models.Loc, to models.Loc) (models.RouteInfo, bool) {
	return c.cachedCalculateRoute(from, to)
}

// NewCachedGeoQuery returns a new cached version of geoQuery
func NewCachedGeoQuery(geoQuery GeoQuery) CachedGeoQuery {
	return &cachedGeoQueryImpl{
		geoQuery:             geoQuery,
		cachedForwardGeocode: cachedForwardGeocodeClosure(geoQuery),
		cachedCalculateRoute: cachedCalculateRouteClosure(geoQuery)}
}

// cachedForwardGeocodeClosure returns a cached version of geoQuery's ForwardGeocode
func cachedForwardGeocodeClosure(geoQuery GeoQuery) (f func(string) (models.Loc, bool)) {
	// The cache
	var addressMap = map[string]models.Loc{}

	f = func(searchString string) (models.Loc, bool) {
		loc, ok := addressMap[searchString]
		if ok {
			return loc, ok
		}
		loc = geoQuery.ForwardGeocode(searchString)
		addressMap[searchString] = loc
		return loc, ok
	}
	return
}

// cachedCalculateRouteClosure returns a cached version of geoQuery's CalculateRoute
func cachedCalculateRouteClosure(geoQuery GeoQuery) (f func(models.Loc, models.Loc) (models.RouteInfo, bool)) {
	// The cache
	var routeInfoMap = map[string]models.RouteInfo{}

	f = func(from models.Loc, to models.Loc) (models.RouteInfo, bool) {
		route := from.Addr + " -> " + to.Addr
		routeInfo, ok := routeInfoMap[route]
		if ok {
			return routeInfo, ok
		}
		routeInfo = geoQuery.CalculateRoute(from, to)
		routeInfoMap[route] = routeInfo
		return routeInfo, ok
	}
	return
}
