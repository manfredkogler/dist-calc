package controllers

import (
	"encoding/csv"
	"os"
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
		"start [km]", "end [km]", "distance [km]", "calculated", "specified", "comment", "route"})
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
