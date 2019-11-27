package main

import (
	"dist-calc/controllers"
	requests "dist-calc/requests"

	distancematrixai "dist-calc/requests/distancematrixai"
	here "dist-calc/requests/here"
	"flag"
)

func main() {
	inFilepathPtr := flag.String("infile", "addresses.csv", "path to a text file containing a list of location specifications/search strings, each given on a new line; this effectively conforms to the csv (comma separated values) format with a single column")
	outFilepathPtr := flag.String("outfile", "results.csv", "path to the generated csv file containing detailed address info, distance and driving time when starting from the preceding location, and other stuff")
	flagStartPoint := flag.Float64("startpoint", 0.0, "starting point in km (use the dot . as comma separator)")
	servicePtr := flag.String("service", "h", "maps service to use (h ... here maps service, d ... distance matrix ai service; default is h)")
	noFileCachePtr := flag.Bool("nofilecache", false, "disable file cache; no cached file is loaded nor stored")
	flag.Parse()

	// Default geo query is the here service
	var geoQuery requests.GeoQuery = here.HereGeoQuery{}
	switch *servicePtr {
	case "d":
		geoQuery = distancematrixai.DistancematrixaiGeoQuery{}
	}

	geoProcessor := controllers.NewProcessor(requests.NewCachedGeoQuery(geoQuery))
	geoProcessor.Start(*inFilepathPtr, *outFilepathPtr, *flagStartPoint, !*noFileCachePtr)
}
