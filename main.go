// It gets config file as the CL argument

package main

import (
	"io/fs"
	"log"
	"os"
	"sort"
	"strconv"
)

var config []string

func main() {

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
		log.Fatal("Cannot read target directory: ", workpath, "\n", err.Error())
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

func filterfiles(files []fs.DirEntry, n int) []fs.DirEntry { // returns cut slice of files
	if n <= len(files)-1 {
		return files[n:]
	}
	return *new([]fs.DirEntry)
}

func processdirs() { // На основании конфига запускает операции
	var n int = len(config) / 2
	for i := 0; i < n; i++ {
		path := config[i*2]
		files := getfiles(path)
		sortfiles(&files)
		filestodelete := filterfiles(files, func() int { i, _ := strconv.Atoi(config[i*2+1]); return i }())
		processfiles(path, filestodelete)
	}
}

func processfiles(path string, files []fs.DirEntry) { // Вызывается из processdirs для выполнения операций
	log.Println("Will be deleted:")
	for _, file := range files {
		log.Println(path + file.Name()) // Will stay anyway for logging purpose
		if !debug {
			err := os.Remove(path + file.Name())
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func remove(s *[]fs.DirEntry, i int) {
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}
