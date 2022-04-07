package main

import (
	"log"
)

func main() {
	initialize()

	log.Println("Starting.")
	if debug {
		log.Println("Debug mode, no file operations will be performed!") // We print this in debug mode
		for _, w := range config {
			log.Println("path: ", w.path)
			log.Println("operations: ", func() []string {
				var res []string
				for _, o := range w.operations {
					o = "\"" + o + "\""
					res = append(res, o)
				}
				return res
			}())
		}
	}
	processjobs()

	log.Println("Exiting normally.") // debug
}
