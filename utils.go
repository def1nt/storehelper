package main

import (
	"os"
	"strings"
)

func readfile(path string) ([]string, error) {
	var raw []byte
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file_t := strings.Split(strings.ReplaceAll(string(raw), "\r", ""), "\n")
	for i := 0; i < len(file_t); i++ {
		file_t[i] = strings.TrimSpace(file_t[i])
		if len(file_t[i]) == 0 {
			removefromslice(&file_t, i)
		}
	}
	return file_t, nil
}

func removefromslice[T any](s *[]T, i int) { // modify given slice by removing item
	*s = append((*s)[:i], (*s)[i+1:]...)
}
