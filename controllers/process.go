package controllers

import (
	"dist-calc/models"
	"dist-calc/requests"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/dimchansky/utfbom"
)

type outFiles struct {
	main      *os.File
	addresses *os.File
	distances *os.File
}

type csvWriters struct {
	main      *csv.Writer
	addresses *csv.Writer
	distances *csv.Writer
}

func createOutFiles(main string, addresses string, distances string) outFiles {
	return outFiles{
		main:      createOutFile(main),
		addresses: createOutFile(addresses),
		distances: createOutFile(distances)}
}

func createCsvWriters(of outFiles) csvWriters {
	return csvWriters{
		main:      createCsvWriter(of.main),
		addresses: createCsvWriter(of.addresses),
		distances: createCsvWriter(of.distances)}
}

func (o *outFiles) close() {
	checkedClose(o.main)
	checkedClose(o.addresses)
	checkedClose(o.distances)
}

func checkedClose(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func createOutFile(filepath string) *os.File {
	outFile, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	// Write the UTF-8 BOM header for Excel to open it with correct encoding
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	_, err = outFile.Write(bomUtf8)
	if err != nil {
		panic(err)
	}

	return outFile
}

func createCsvWriter(outFile *os.File) *csv.Writer {
	w := csv.NewWriter(outFile)
	w.Comma = ';'
	return w
}

func (w *csvWriters) writeColumnHeaders() {
	// Write column headers
	checkedWrite(w.main, []string{
		"start [km]", "end [km]", "distance [km]", "route"})
	checkedWrite(w.addresses, []string{
		"address specified", "address found", "latitude", "longitude"})
	checkedWrite(w.distances, []string{
		"route", "distance [m]", "travel time [s]"})
}

func checkedWrite(w *csv.Writer, record []string) {
	err := w.Write(record)
	if err != nil {
		panic(err)
	}
}

func (w *csvWriters) flush() {
	checkedFlush(w.main)
	checkedFlush(w.addresses)
	checkedFlush(w.distances)
}

func checkedFlush(w *csv.Writer) {
	w.Flush()
	err := w.Error()
	if err != nil {
		panic(err)
	}
}

// ProcessAdressList traverses the address list and generates the output files
func ProcessAdressList(inFilepath string, outFilepath string, startPoint float64) {
	inFile, err := os.Open(inFilepath)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// Skip the BOM if there, otherwise it will be part of the first address parsed
	bomSkippedInFile := utfbom.SkipOnly(inFile)

	r := csv.NewReader(bomSkippedInFile)
	r.Comma = ';'
	r.TrimLeadingSpace = true
	// r.ReuseRecord = true

	outFilePathWoExt := strings.TrimSuffix(outFilepath, path.Ext(outFilepath))
	outFiles := createOutFiles(outFilepath, outFilePathWoExt+"-addresses-retrieved.csv", outFilePathWoExt+"-distances-retrieved.csv")
	defer outFiles.close()

	csvWriters := createCsvWriters(outFiles)
	csvWriters.writeColumnHeaders()

	startKm := startPoint
	endKm := startKm

	fromSpec, _ := readNextAddressSpec(r, csvWriters.main, &startKm, &endKm)
	if fromSpec == "" {
		return
	}
	from := handleForwardGeocode(fromSpec, csvWriters.addresses)

	for {
		fmt.Println("----------------------------------------------------------------------")

		// var toSpec string
		toSpec, distanceSpecHandled := readNextAddressSpec(r, csvWriters.main, &startKm, &endKm)
		if toSpec == "" {
			break
		}
		to := handleForwardGeocode(toSpec, csvWriters.addresses)

		// If we have handled a distance specification we cannot yet calculate a distance and
		// we need safe the to as the from and read in the next address spec by restarting this loop
		if distanceSpecHandled {
			fromSpec = toSpec
			from = to
			continue
		}

		routeInfo := cachedCalculateRoute(from, to)
		fmt.Println("RouteInfo: ", routeInfo)

		distanceKm := float64(routeInfo.Distance) / 1000
		startKm = endKm
		endKm += distanceKm

		// Write next line / location
		routeSpec := fromSpec + " -> " + toSpec
		checkedWrite(csvWriters.main, []string{
			floatToString(startKm), floatToString(endKm), floatToString(distanceKm), routeSpec})
		checkedWrite(csvWriters.distances, []string{
			routeSpec, strconv.FormatInt(routeInfo.Distance, 10), strconv.FormatInt(routeInfo.TravelTime, 10)})

		fromSpec = toSpec
		from = to
	}

	csvWriters.flush()
}

var cachedForwardGeocode = requests.CachedForwardGeocodeClosure()

func handleForwardGeocode(addrSpec string, addresses *csv.Writer) models.Loc {
	loc, fromCache := cachedForwardGeocode(addrSpec)
	// Write a new record to the address file
	if !fromCache {
		checkedWrite(addresses, []string{addrSpec, loc.Addr, loc.Lat, loc.Lng})
	}
	return loc
}

var cachedCalculateRoute = requests.CachedCalculateRouteClosure()

// readNextAddressSpec reads the next address spec; any distance specification before is properly handled
// return value:
//   string:	address specification (empty string if eof)
//   bool:		any distance specification handled
func readNextAddressSpec(r *csv.Reader, w *csv.Writer, startKm *float64, endKm *float64) (string, bool) {
	var distanceSpecHandled bool

	record, err := r.Read()
	if err == io.EOF {
		return "", distanceSpecHandled
	}
	if err != nil {
		panic(err)
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
		checkedWrite(w, []string{
			floatToString(*startKm), floatToString(*endKm), floatToString(distanceSpecKm), "-"})

		record, err := r.Read()
		if err == io.EOF {
			// Signaling eof
			fromSpec = ""
			break
		}
		if err != nil {
			panic(err)
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
