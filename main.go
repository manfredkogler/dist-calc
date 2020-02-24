package main

import (
	"flag"

	"github.com/manfredkogler/dist-calc/controllers"
	"github.com/manfredkogler/dist-calc/requests"

	distancematrixai "github.com/manfredkogler/dist-calc/requests/distancematrixai"
	here "github.com/manfredkogler/dist-calc/requests/here"
)

func main() {
	inFilepathPtr := flag.String("infile", "addresses.csv", "path to a text file containing a list of location specifications/search strings")
	outFilepathPtr := flag.String("outfile", "results.csv", "path to the generated main results csv file")
	flagStartPoint := flag.Float64("startpoint", 0.0, "starting point in km (default 0.0)")
	servicePtr := flag.String("service", "h", "maps service to use (h - here maps service, d - distance matrix ai service)")
	noFileCachePtr := flag.Bool("nofilecache", false, "disable file cache; nothing is loaded from or stored to a local file")
	spreadBasePtr := flag.Float64("spreadbase", 0.3, "spread base in km")
	spreadFactorPtr := flag.Float64("spreadfactor", 0.005, "spread factor")
	flag.Parse()

	// Set the requested geo query service
	var geoQuery requests.GeoQuery
	switch *servicePtr {
	case "d":
		geoQuery = distancematrixai.DistancematrixaiGeoQuery{}
	default:
		// Default geo query is the here service
		geoQuery = here.HereGeoQuery{}
	}

	// We use the cached version of the geo query
	geoProcessor := controllers.NewProcessor(requests.NewCachedGeoQuery(geoQuery), *spreadBasePtr, *spreadFactorPtr)
	geoProcessor.Start(*inFilepathPtr, *outFilepathPtr, *flagStartPoint, !*noFileCachePtr)
}
