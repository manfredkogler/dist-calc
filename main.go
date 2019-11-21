package main

import (
	"dist-calc/models"
	"dist-calc/requests"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	// in := `Lambertgasse 4/2/3, Wien
	// Schottenring 1, Wien
	// Rückertgasse 22/20, Wien
	// Freyung 1, Wien
	// Am Waldrand 1, Adlwang`

	inFilepathPtr := flag.String("infile", "addresses.csv", "path to a text file containing a list of location specifications/search strings, each given on a new line; this effectively conforms to the csv (comma separated values) format with a single column")
	outFilepathPtr := flag.String("outfile", "results.csv", "path to the generated csv file containing detailed address info, distance and driving time when starting from the preceding location, and other stuff")
	// outputLatLng := flag.Bool("lat-long", false, "also print latitude and longitude")
	flagStartPoint := flag.Float64("startpoint", 0.0, "starting point in km (use the dot . as comma separator)")
	flagStartDistance := flag.Float64("startdistance", 0.0, "starting distance in km (use the dot . as comma separator)")
	startRoute := flag.String("startroute", "-", "starting route description")

	flag.Parse()

	// Open the input file
	inFile, err := os.Open(*inFilepathPtr)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	r := csv.NewReader(inFile)

	// r := csv.NewReader(strings.NewReader(in))
	r.Comma = ';'
	r.TrimLeadingSpace = true
	// r.ReuseRecord = true

	outFile, err := os.Create(*outFilepathPtr)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Write the UTF-8 BOM header for Excel to open it with correct encoding
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	outFile.Write(bomUtf8)

	csvWriter := csv.NewWriter(outFile)
	csvWriter.Comma = ';'

	// Write column headers
	err = csvWriter.Write([]string{
		"address specified", "address found", "latitude", "longitude",
		"distance [m]", "travel time [s]",
		"start [km]", "end [km]", "distance [km]",
		"route"})

	record, err := r.Read()
	if err == io.EOF {
		log.Fatal(err)
	}

	fromSpec := record[0]
	from := requests.ForwardGeocode(fromSpec)

	distanceKm := *flagStartDistance
	startKm := *flagStartPoint
	endKm := startKm + distanceKm

	// Write first line / starting location
	err = csvWriter.Write([]string{
		fromSpec, from.Addr, from.Lat, from.Lng,
		"0", "0",
		floatToString(startKm), floatToString(endKm), floatToString(distanceKm),
		*startRoute})

	var to models.Loc
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("----------------------------------------------------------------------")
		fmt.Println(record)

		toSpec := record[0]
		to = requests.ForwardGeocode(toSpec)

		routeInfo := requests.CalculateRoute(from, to)
		fmt.Println("RouteInfo: ", routeInfo)

		distanceKm = float64(routeInfo.Distance) / 1000
		startKm = endKm
		endKm += distanceKm

		// Write next line / location
		err = csvWriter.Write([]string{
			toSpec, to.Addr, to.Lat, to.Lng,
			strconv.FormatInt(routeInfo.Distance, 10), strconv.FormatInt(routeInfo.TravelTime, 10),
			floatToString(startKm), floatToString(endKm), floatToString(distanceKm),
			fromSpec + " -> " + toSpec})

		fromSpec = toSpec
		from = to
	}

	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		// an error occurred during the flush
	}

	// routeInfo := requests.CalculateRoute(requests.Locations[requests.Lambertgasse], requests.Locations[requests.Schottenring])
	// routeInfo = requests.CalculateRoute(requests.Locations[requests.Schottenring], requests.Locations[requests.Lambertgasse])
	// _ = routeInfo

	// location := requests.ForwardGeocode("Schottenring 1 Wien")
	// location = requests.ForwardGeocode("Lambertgasse 4 Wien")
	// _ = location
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	value := strconv.FormatFloat(inputNum, 'f', 1, 64)
	// Use comma instead of dot as decimal "point" for Excel to properly handle it
	return strings.Replace(value, ".", ",", -1)
}
