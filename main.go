package main

import (
	"io/fs"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	log.Println("Starting.")

	initialize()
	if debug {
		log.Println("Debug mode, no file operations will be performed!")
	}
	processdirs()

	log.Println("Exiting normally.") // debug
}

func getfiles(workpath string) []fs.DirEntry {
	files, err := os.ReadDir(workpath)
	if err != nil {
		log.Println("Cannot read target directory: ", workpath, "\n", err.Error()) // not fatal, just log and skip
		return nil
	}
	for i := 0; i < len(files); i++ { // Safe for work
		if files[i].IsDir() {
			remove(&files, i)
			i--
		}
	}
	return files
}

func sortfiles(files *[]fs.DirEntry) *[]fs.DirEntry { // Sorts original list in place by date \
	sort.Slice(*files, func(i, j int) bool {
		fi, _ := (*files)[i].Info()
		fj, _ := (*files)[j].Info()
		return fi.ModTime().Unix() > fj.ModTime().Unix() // Sort high to low, to cut top N files
	})
	return files
}

func filterfiles(files []fs.DirEntry, ops string) []fs.DirEntry { // returns cut slice of files
	var s strings.Builder
	var t string
	var n int
	for _, r := range ops {
		if r == ':' {
			t = s.String()
			s.Reset()
			continue
		}
		s.WriteRune(r)
	}
	n, err := strconv.Atoi(s.String())
	if err != nil {
		log.Println(err.Error())
	}

	switch t {
	case "n":
		if n < 1 {
			break // and return empty
		}
		if n >= len(files) {
			return files
		}
		if n < len(files) {
			return files[:n]
		}
	case "o":
		if n < 1 {
			break
		}
		if n >= len(files) {
			return files
		}
		if n < len(files) {
			return files[len(files)-n:]
		}
	case "y":
	case "e":
	}

	return *new([]fs.DirEntry)
}

func processdirs() { // На основании конфига запускает операции
	for _, wi := range config {
		path := wi.path
		files := getfiles(path)
		if files == nil { // if cannot read directory, skip this one, error logged in getfiles()
			continue
		}
		sortfiles(&files)
		filtered := filterfiles(files, wi.operations[0])
		processfiles(path, files, filtered, wi.operations[1])
	}
}

func processfiles(path string, files []fs.DirEntry, filtered []fs.DirEntry, operation string) { // Вызывается для выполнения операций
	log.Println("Working in: " + path)
	switch operation {
	case "d":
		log.Println("Will be deleted:")
		for _, file := range filtered {
			log.Println(path + file.Name()) // Will stay anyway for logging purpose
			if !debug {
				err := os.Remove(path + file.Name())
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	case "k":
		log.Println("Will be deleted:")
	deleting:
		for _, file := range files {
			for _, keepfile := range filtered {
				if keepfile.Name() == file.Name() {
					continue deleting
				}
			}
			log.Println(path + file.Name()) // Will stay anyway for logging purpose
			if !debug {
				err := os.Remove(path + file.Name())
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

func remove(s *[]fs.DirEntry, i int) { // modify given slice by removing item
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}
