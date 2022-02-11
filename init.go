package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

var debug bool = false // Отладочный режим, не выполнять операции, только выводит отчёт, ключ -t

type workitem struct {
	path       string
	operations []string
}

var config1 []workitem

func initialize() error {
	if len(os.Args) == 1 {
		log.Fatal("No config file specified as a first argument!")
	}

	if len(os.Args) > 2 {
		for _, s := range os.Args {
			switch s {
			case "t", "-t":
				debug = true
			}
		}
	}

	readconfig()

	return nil
}

func readconfig() {
	var conffile string = os.Args[1]
	conffile = path.Clean(conffile)
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

	readworkitems()
	if debug { // Probably will go elswhere, we print this in debug mode
		for _, w := range config1 {
			fmt.Println("path: ", w.path)
			fmt.Println("operations: ", func() []string {
				var res []string
				for _, o := range w.operations {
					o = "\"" + o + "\""
					res = append(res, o)
				}
				return res
			}())
		}
	}

	for i := range config { // Every 2nd line is a path and needs to be cleaned
		if i%2 == 0 {
			config[i] = path.Clean(config[i])
		}
	}
}

func readworkitems() {
	for i := 0; i+1 < len(config); i += 2 {
		var wi workitem
		wi.path = path.Clean(config[i])
		wi.operations = processoperations(config[i+1])
		config1 = append(config1, wi)
	}
}

func processoperations(operations string) []string {
	var s strings.Builder
	var ops []string
	for _, r := range operations {
		if r == ' ' {
			if s.Len() != 0 {
				ops = append(ops, s.String())
			}
			s.Reset()
			continue
		}
		s.WriteRune(r)
	}
	if s.Len() != 0 { // If the string is NOT closed by single space
		ops = append(ops, s.String())
	}

	return ops
}
