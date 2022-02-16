package utils

import (
	"encoding/csv"
	"os"
)

// CsvReadAll reads a file and returns contents as a [][]string
func CsvReadAll(filepath string, ignoreHeader bool) ([][]string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0555)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if ignoreHeader {
		return data[1:], nil
	}
	return data, nil
}
