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
	var file_t []string
	var param strings.Builder

	for _, b := range raw {
		if b == 13 {
			continue
		}
		if b == 10 || b == 0 {
			if param.Len() > 0 {
				file_t = append(file_t, param.String())
				param.Reset()
			}
			continue
		}
		param.WriteByte(b)
	}
	if param.Len() > 0 {
		file_t = append(file_t, param.String())
	}
	return file_t, nil
}

func removefromslice[T any](s *[]T, i int) { // modify given slice by removing item
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}
