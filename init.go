package main

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
)

var debug bool = false // Отладочный режим, не выполнять операции, только выводит отчёт, ключ -t

func initialize() error {
	if len(os.Args) == 1 {
		log.Fatal("No config file specified as a first argument")
	}

	if len(os.Args) > 2 {
		for _, s := range os.Args {
			switch s {
			case "t", "-t":
				debug = true
			}
		}
	}

	var conffile string = os.Args[1]
	path.Clean(conffile)
	var bytes []byte
	var err error
	bytes, err = os.ReadFile(conffile)

	if err != nil && err != io.EOF {
		log.Fatal("Cannot read config, error: ", err.Error())
	}

	var param strings.Builder
	for _, b := range bytes {
		if b == 35 { // # — стоп символ, пока и так сойдёт
			break
		}
		if b == 13 {
			continue
		}
		if b == 10 || b == 0 {
			if param.Len() > 0 {
				config = append(config, param.String())
				param.Reset()
			}
			continue
		}
		param.WriteByte(b)
	}
	if param.Len() > 0 {
		config = append(config, param.String())
	}

	for i := range config { // Every 2nd line is a path and needs to be cleaned
		if i%2 == 0 {
			config[i] = path.Clean(config[i])
		}

	}

	return nil
}
