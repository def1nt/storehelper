package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func processjobs() {
	for _, wi := range config {
		dooperations(wi)
	}
}

func getfiles(workpath string, ftype string) []fs.DirEntry {
	files, err := os.ReadDir(workpath)
	if err != nil {
		log.Println("Cannot read target directory:", err.Error()) // not fatal, just log and skip
		return nil
	}
	for i := 0; i < len(files); i++ { // Safe for work: len evaluates every loop
		if (ftype == "f" || ftype == "files") && files[i].IsDir() {
			removefromslice(&files, i)
			i--
		}
		if (ftype == "d" || ftype == "dirs") && !files[i].IsDir() {
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
	case "n", "newest":
		if n < 1 {
			break // and return empty implicitly
		}
		if n >= len(files) {
			return files
		}
		if n < len(files) {
			return files[:n]
		}
	case "o", "oldest":
		if n < 1 {
			break
		}
		if n >= len(files) {
			return files
		}
		if n < len(files) {
			return files[len(files)-n:]
		}
	case "y", "younger":
		var t = make([]fs.DirEntry, 0, len(files))
		var age float64 = float64(n * 24)
		for _, f := range files {
			fi, _ := f.Info()
			if time.Since(fi.ModTime()).Hours() <= age {
				t = append(t, f)
			}
		}
		return t
	case "e", "elder":
		var t = make([]fs.DirEntry, 0, len(files))
		var age float64 = float64(n * 24)
		for _, f := range files {
			fi, _ := f.Info()
			if time.Since(fi.ModTime()).Hours() > age {
				t = append(t, f)
			}
		}
		return t
	case "w", "weekday":
		var t = make([]fs.DirEntry, 0, cap(files))
		for _, f := range files {
			fi, _ := f.Info()
			wd := int(fi.ModTime().Weekday())
			if wd == 0 {
				wd = 7
			}
			if strings.ContainsAny(param, strconv.Itoa(wd)) {
				t = append(t, f)
			}
		}
		return t
	case "m", "monthday":
		var t = make([]fs.DirEntry, 0, cap(files))
		for _, f := range files {
			fi, _ := f.Info()
			if fi.ModTime().Day() == n {
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
		case "d", "f", "df",
			"dirs", "files":
			selectiontype = op
			files = getfiles(path, op)
			if files == nil {
				return // if nil then there was an error reading directory, skipping workitem
			}
		case "n", "o", "y", "e", "w", "m",
			"newest", "oldest", "younger", "elder", "weekday", "monthday":
			files = filter(op, param, files)
		case "r", "remove":
			remove(path, files)
		case "k", "keep":
			remove(path, invert(path, selectiontype, files))
		case "c", "copy":
			copy(path, param, files)
		case "v", "verify":
			verify(files, param, path)
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
			err := os.RemoveAll(path + file.Name())
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
		dfile.Close()
		sstat, _ := sfile.Stat()
		sfile.Close()
		err = os.Chtimes(destination+f.Name(), sstat.ModTime(), sstat.ModTime())
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("Copied", n, "bytes from", source+sfile.Name(), "to", destination)
	}
}
