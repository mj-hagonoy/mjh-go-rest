package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// CsvReadAll reads a file and returns contents as a [][]string
func CsvReadAll(filepath string, ignoreHeader bool) ([][]string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0555)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return CsvRead(file, ignoreHeader)
}

func CsvRead(r io.Reader, ignoreHeader bool) ([][]string, error) {
	reader := csv.NewReader(r)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CsvRead: %v", err)
	}
	if ignoreHeader {
		return data[1:], nil
	}
	return data, nil
}
