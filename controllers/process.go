package controllers

import (
	"dist-calc/common"
	"dist-calc/models"
	"dist-calc/requests"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dimchansky/utfbom"
)

// NewProcessor returns a new processor for some maps service implementing the CachedGeoQuery interface
func NewProcessor(cachedGeoQuery requests.CachedGeoQuery, spreadBase float64, spreadFactor float64) *Processor {
	return &Processor{
		cachedGeoQuery:        cachedGeoQuery,
		addressRecordWritten:  map[string]bool{},
		distanceRecordWritten: map[string]bool{},
		spreadBase:            spreadBase,
		spreadFactor:          spreadFactor}
}

// Processor processes geo requests
type Processor struct {
	cachedGeoQuery        requests.CachedGeoQuery
	addressRecordWritten  map[string]bool
	distanceRecordWritten map[string]bool
	spreadBase            float64
	spreadFactor          float64
}

// Start loads caches, processes the address list and stores the caches
func (p Processor) Start(inFilepath string, outFilepath string, startPoint float64, useFileCache bool) {
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

	if useFileCache {
		p.cachedGeoQuery.LoadCaches(csvWriters.addresses, csvWriters.distances)
	}

	// Random package needed for adding spread to distance calculations
	rand.Seed(time.Now().UnixNano())

	p.ProcessAdressList(r, csvWriters, startPoint)
	if useFileCache {
		p.cachedGeoQuery.StoreCaches()
	}
}

// ProcessAdressList traverses the address list and generates the output files
func (p Processor) ProcessAdressList(r *csv.Reader, csvWriters csvWriters, startPoint float64) {
	startKm := startPoint
	endKm := startKm

	fromSpec, _ := readNextAddressSpec(r, csvWriters.main, &startKm, &endKm)
	if fromSpec == "" {
		return
	}
	from := p.handleForwardGeocode(fromSpec, csvWriters.addresses)

	for {
		fmt.Println("----------------------------------------------------------------------")

		// var toSpec string
		toSpec, distanceSpecHandled := readNextAddressSpec(r, csvWriters.main, &startKm, &endKm)
		if toSpec == "" {
			break
		}
		to := p.handleForwardGeocode(toSpec, csvWriters.addresses)

		// If we have handled a distance specification we cannot yet calculate a distance and
		// we need safe the to as the from and read in the next address spec by restarting this loop
		if distanceSpecHandled {
			fromSpec = toSpec
			from = to
			continue
		}

		routeInfo, _ := p.cachedGeoQuery.CalculateRoute(from, to)
		fmt.Println("RouteInfo: ", routeInfo)

		distanceKm := float64(routeInfo.Distance) / 1000

		// Add real world spread to calculated distance
		maxSpreadToAdd := p.spreadBase + distanceKm*p.spreadFactor
		distanceKmSpread := distanceKm + rand.Float64()*maxSpreadToAdd

		startKm = endKm
		endKm += distanceKmSpread

		// Write next line / location
		routeSpec := fromSpec + " -> " + toSpec
		checkedWrite(csvWriters.main, []string{
			floatToString(startKm), floatToString(endKm), floatToString(distanceKmSpread), "x", "", "", routeSpec})

		// Write a new record to the distance file
		if !p.distanceRecordWritten[routeSpec] {
			checkedWrite(csvWriters.distances, []string{
				routeSpec,
				strconv.FormatInt(int64(routeInfo.Distance), 10),
				floatToString(distanceKm),
				strconv.FormatInt(int64(routeInfo.TravelTime), 10)})
			p.distanceRecordWritten[routeSpec] = true
		}

		fromSpec = toSpec
		from = to
	}

	csvWriters.flush()
}

func (p Processor) handleForwardGeocode(addrSpec string, addresses *csv.Writer) models.Loc {
	loc, _ := p.cachedGeoQuery.ForwardGeocode(addrSpec)
	// Write a new record to the address file
	if !p.addressRecordWritten[addrSpec] {
		checkedWrite(addresses, []string{addrSpec, loc.Addr, loc.Lat, loc.Lng})
		p.addressRecordWritten[addrSpec] = true
	}
	return loc
}

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
			floatToString(*startKm), floatToString(*endKm), floatToString(distanceSpecKm), "", "x", "", "-"})

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
	return common.FloatToString(inputNum, 1)
}
