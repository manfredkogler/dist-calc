package requests

import (
	"dist-calc/models"
	"encoding/csv"
	"encoding/gob"
	"os"
	"strconv"
)

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
	// Load the cache from a file
	LoadCaches(addresses *csv.Writer, distances *csv.Writer)
	// Store the cache to a file
	StoreCaches()
	// ForwardGeocode is a cached version of GeoQuery.ForwardGeocode. The bool returned signals if the returned entry is taken from the cache.
	ForwardGeocode(searchString string) (models.Loc, bool)
	// CalculateRoute is a cached version of GeoQuery.CalculateRoute. The bool returned signals if the returned entry is taken from the cache.
	CalculateRoute(from models.Loc, to models.Loc) (models.RouteInfo, bool)
}

// cachedGeoQueryImpl implements a cached version of GeoQuery
type cachedGeoQueryImpl struct {
	geoQuery GeoQuery
	// Address cache
	addressMap map[string]models.Loc
	// Route info cache
	routeInfoMap map[string]models.RouteInfo
}

// Load the caches from stored files
func (c cachedGeoQueryImpl) LoadCaches(addresses *csv.Writer, distances *csv.Writer) {
	c.loadAddressCache(addresses)
	c.loadRouteInfoCache(distances)
}

// Load the address cache from a file
func (c cachedGeoQueryImpl) loadAddressCache(addresses *csv.Writer) {
	decodeFile, err := os.Open("addressMap.gob")
	if err != nil {
		// No file, no cache yet, no issue...
		return
	}
	defer decodeFile.Close()

	decoder := gob.NewDecoder(decodeFile)
	decoder.Decode(&c.addressMap)

	// Write the loaded/cached records to the addresses file
	for k, v := range c.addressMap {
		// Write a new record to the address file
		checkedWrite(addresses, []string{k, v.Addr, v.Lat, v.Lng})
	}
}

func checkedWrite(w *csv.Writer, record []string) {
	err := w.Write(record)
	if err != nil {
		panic(err)
	}
}

// Load the route info cache from a file
func (c cachedGeoQueryImpl) loadRouteInfoCache(distances *csv.Writer) {
	decodeFile, err := os.Open("routeInfoMap.gob")
	if err != nil {
		// No file, no cache yet, no issue...
		return
	}
	defer decodeFile.Close()

	decoder := gob.NewDecoder(decodeFile)
	decoder.Decode(&c.routeInfoMap)

	// Write the loaded/cached records to the distances file
	for k, v := range c.routeInfoMap {
		// Write a new record to the distance file
		checkedWrite(distances, []string{
			k, strconv.FormatInt(int64(v.Distance), 10), strconv.FormatInt(int64(v.TravelTime), 10)})
	}
}

// Store the caches to files
func (c cachedGeoQueryImpl) StoreCaches() {
	c.storeAddressCache()
	c.storeRouteInfoCache()
}

// Store the address cache to a file
func (c cachedGeoQueryImpl) storeAddressCache() {
	encodeFile, err := os.Create("addressMap.gob")
	if err != nil {
		panic(err)
	}

	encoder := gob.NewEncoder(encodeFile)
	if err := encoder.Encode(c.addressMap); err != nil {
		panic(err)
	}
	encodeFile.Close()
}

// Store the route info cache to a file
func (c cachedGeoQueryImpl) storeRouteInfoCache() {
	encodeFile, err := os.Create("routeInfoMap.gob")
	if err != nil {
		panic(err)
	}

	encoder := gob.NewEncoder(encodeFile)
	if err := encoder.Encode(c.routeInfoMap); err != nil {
		panic(err)
	}
	encodeFile.Close()
}

// CachedForwardGeocode implements a cached version of GeoQuery.ForwardGeocode
func (c cachedGeoQueryImpl) ForwardGeocode(searchString string) (models.Loc, bool) {
	loc, ok := c.addressMap[searchString]
	if ok {
		return loc, ok
	}
	loc = c.geoQuery.ForwardGeocode(searchString)
	c.addressMap[searchString] = loc
	return loc, ok
}

// CachedCalculateRoute implements a cached version of GeoQuery.CalculateRoute
func (c cachedGeoQueryImpl) CalculateRoute(from models.Loc, to models.Loc) (models.RouteInfo, bool) {
	route := from.Addr + " -> " + to.Addr
	routeInfo, ok := c.routeInfoMap[route]
	if ok {
		return routeInfo, ok
	}
	routeInfo = c.geoQuery.CalculateRoute(from, to)
	c.routeInfoMap[route] = routeInfo
	return routeInfo, ok
}

// NewCachedGeoQuery returns a new cached version of geoQuery
func NewCachedGeoQuery(geoQuery GeoQuery) CachedGeoQuery {
	return &cachedGeoQueryImpl{
		geoQuery:     geoQuery,
		addressMap:   map[string]models.Loc{},
		routeInfoMap: map[string]models.RouteInfo{},
	}
}
