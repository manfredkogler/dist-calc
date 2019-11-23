package controllers

import (
	"dist-calc/models"
	"dist-calc/requests"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func ProcessAdressList(inFilepath string, outFilepath string, startPoint float64) {
	inFile, err := os.Open(inFilepath)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	r := csv.NewReader(inFile)
	r.Comma = ';'
	r.TrimLeadingSpace = true
	// r.ReuseRecord = true

	outFile, err := os.Create(outFilepath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Write the UTF-8 BOM header for Excel to open it with correct encoding
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	outFile.Write(bomUtf8)

	w := csv.NewWriter(outFile)
	w.Comma = ';'

	// Write column headers
	err = w.Write([]string{
		"address specified", "address found", "latitude", "longitude",
		"distance [m]", "travel time [s]",
		"start [km]", "end [km]", "distance [km]",
		"route"})

	startKm := startPoint
	endKm := startKm

	fromSpec, distanceSpecHandled := readNextAddressSpec(r, w, &startKm, &endKm)
	if fromSpec == "" {
		return
	}
	from := requests.ForwardGeocode(fromSpec)

	var to models.Loc
	for {
		fmt.Println("----------------------------------------------------------------------")

		var toSpec string
		toSpec, distanceSpecHandled = readNextAddressSpec(r, w, &startKm, &endKm)
		if toSpec == "" {
			break
		}
		to = requests.ForwardGeocode(toSpec)

		// If we have handled a distance specification we cannot yet calculate a distance and
		// we need safe the to as the from and read in the next address spec by restarting this loop
		if distanceSpecHandled {
			fromSpec = toSpec
			from = to
			continue
		}

		routeInfo := requests.CalculateRoute(from, to)
		fmt.Println("RouteInfo: ", routeInfo)

		distanceKm := float64(routeInfo.Distance) / 1000
		startKm = endKm
		endKm += distanceKm

		// Write next line / location
		err = w.Write([]string{
			toSpec, to.Addr, to.Lat, to.Lng,
			strconv.FormatInt(routeInfo.Distance, 10), strconv.FormatInt(routeInfo.TravelTime, 10),
			floatToString(startKm), floatToString(endKm), floatToString(distanceKm),
			fromSpec + " -> " + toSpec})

		fromSpec = toSpec
		from = to
	}

	w.Flush()
	err = w.Error()
	if err != nil {
		// an error occurred during the flush
		panic(err)
	}
}

// readNextAddressSpec reads the next address spec; any distance specifications before are handled before
// return value:
//   string:	address specification (empty of eof)
//   bool:		distance specification(s) handled
func readNextAddressSpec(r *csv.Reader, w *csv.Writer, startKm *float64, endKm *float64) (string, bool) {
	var distanceSpecHandled bool

	record, err := r.Read()
	if err == io.EOF {
		return "", distanceSpecHandled
	}
	if err != nil {
		log.Fatal(err)
	}
	fromSpec := record[0]

	// Check for distance specification, i.e. no address but a distance is given as a float value instead
	for {
		var distanceSpecKm float64
		_, err = fmt.Fscanf(strings.NewReader(fromSpec), "%f", &distanceSpecKm)
		// If no distance is given we have nothing to do
		if err != nil {
			break
		}
		distanceSpecHandled = true

		*startKm = *endKm
		*endKm += distanceSpecKm

		// Write direct distance record
		err = w.Write([]string{
			"-", "-", "-", "-",
			"-", "-",
			floatToString(*startKm), floatToString(*endKm), floatToString(distanceSpecKm),
			"-"})

		record, err := r.Read()
		if err == io.EOF {
			fromSpec = ""
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fromSpec = record[0]
	}

	return fromSpec, distanceSpecHandled
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	value := strconv.FormatFloat(inputNum, 'f', 1, 64)
	// Use comma instead of dot as decimal "point" for Excel to properly handle it
	return strings.Replace(value, ".", ",", -1)
}
