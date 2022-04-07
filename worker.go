package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func processjobs() {
	for _, wi := range config {
		dooperations(wi)
	}
}

func getfiles(workpath string, ftype string) []fs.DirEntry {
	files, err := os.ReadDir(workpath)
	if err != nil {
		log.Println("Cannot read target directory: ", workpath, "\n", err.Error()) // not fatal, just log and skip
		return nil
	}
	for i := 0; i < len(files); i++ { // Safe for work
		if ftype == "f" && files[i].IsDir() {
			removefromslice(&files, i)
			i--
		}
		if ftype == "d" && !files[i].IsDir() {
			removefromslice(&files, i)
			i--
		}
	}
	sortfiles(&files)
	return files
}

func sortfiles(files *[]fs.DirEntry) *[]fs.DirEntry {
	sort.Slice(*files, func(i, j int) bool {
		fi, _ := (*files)[i].Info()
		fj, _ := (*files)[j].Info()
		return fi.ModTime().Unix() > fj.ModTime().Unix() // Sort high to low, to cut top N files
	})
	return files
}

func filter(op string, param string, files []fs.DirEntry) []fs.DirEntry {
	n, err := strconv.Atoi(param)
	if err != nil {
		log.Println(err.Error())
		return make([]fs.DirEntry, 0)
	}

	switch op {
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
	case "w":
		var t = make([]fs.DirEntry, 0, cap(files))
		for _, f := range files {
			fi, _ := f.Info()
			if int(fi.ModTime().Weekday()) == n%7 {
				t = append(t, f)
			}
		}
		return t
	}
	return make([]fs.DirEntry, 0)
}

func invert(path string, ftype string, files []fs.DirEntry) []fs.DirEntry {
	original := getfiles(path, ftype)
	var inverted []fs.DirEntry = make([]fs.DirEntry, 0, cap(original))

matching:
	for _, o := range original {
		for _, t := range files {
			if o.Name() == t.Name() {
				continue matching
			}
		}
		inverted = append(inverted, o)
	}
	return inverted
}

func dooperations(wi workitem) { // We perform these with A SINGLE WORKITEM multiple OPERATIONS
	path := wi.path
	log.Println("Working in:", path)
	var files []fs.DirEntry
	var selectiontype string = "f" // by default we only select files
	for _, ops := range wi.operations {
		op, param := parse(ops)
		// log.Println("op:", op, "param:", param) // Can use for step-by-step debug
		switch op {
		case "d", "f", "df":
			selectiontype = op
			files = getfiles(path, op)
		case "n", "o", "w":
			files = filter(op, param, files)
		case "r":
			remove(path, files)
		case "k":
			remove(path, invert(path, selectiontype, files))
		case "c":
			copy(path, param, files)
		}
	}
}

func parse(ops string) (op string, param string) {
	var s strings.Builder
	var t string
	for _, r := range ops {
		if r == ':' {
			t = s.String()
			s.Reset()
			continue
		}
		s.WriteRune(r)
	}
	if t == "" {
		t = s.String()
		s.Reset()
	}
	return t, s.String()
}

func remove(path string, files []fs.DirEntry) {
	log.Println("Will be deleted:")
	for _, file := range files {
		log.Println(path + file.Name()) // Will stay anyway for logging purpose
		if !debug {
			var err error
			if file.IsDir() {
				err = os.RemoveAll(path + file.Name())
			} else {
				err = os.Remove(path + file.Name()) // Actually, we can use RemoveAll for everything
			}
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func copy(source string, destination string, files []fs.DirEntry) {
	log.Println("Will be copied:")
	for _, f := range files {
		if debug {
			log.Println(source+f.Name(), "to", destination)
			continue
		}

		sfile, err := os.Open(source + f.Name())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		dfile, err := os.Create(destination + f.Name())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		n, err := io.Copy(dfile, sfile)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println("Copied", n, "bytes from", source+sfile.Name(), "to", destination)
		sfile.Close()
		dfile.Close()
	}
}

func removefromslice(s *[]fs.DirEntry, i int) { // modify given slice by removing item
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}
