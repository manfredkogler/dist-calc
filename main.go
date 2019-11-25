package main

import (
	"dist-calc/controllers"
	requests "dist-calc/requests"
	here "dist-calc/requests/here"
	"flag"
)

func main() {
	inFilepathPtr := flag.String("infile", "addresses.csv", "path to a text file containing a list of location specifications/search strings, each given on a new line; this effectively conforms to the csv (comma separated values) format with a single column")
	outFilepathPtr := flag.String("outfile", "results.csv", "path to the generated csv file containing detailed address info, distance and driving time when starting from the preceding location, and other stuff")
	flagStartPoint := flag.Float64("startpoint", 0.0, "starting point in km (use the dot . as comma separator)")
	flag.Parse()

	hereProcessor := controllers.NewProcessor(requests.NewCachedGeoQuery(here.HereGeoQuery{}))
	hereProcessor.ProcessAdressList(*inFilepathPtr, *outFilepathPtr, *flagStartPoint)
}
