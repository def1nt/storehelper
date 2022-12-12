package main

import (
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

var config []workitem

func initialize() {
	if len(os.Args) == 1 {
		log.Fatal("No config file specified as a first argument!")
	}

	if len(os.Args) > 2 {
		for i, s := range os.Args {
			switch s {
			case "t", "-t":
				debug = true
			case "l", "-l":
				file, err := os.OpenFile(os.Args[i+1], os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
				if err != nil {
					log.Fatal(err.Error())
				}
				log.SetOutput(file)
			}
		}
	}

	readconfig()
}

func readconfig() {
	var conffile string = os.Args[1]
	conffile = path.Clean(conffile)

	temp, err := readfile(conffile)
	if err != nil {
		log.Fatal("Cannot read config, error: ", err.Error())
	}

	var config_t []string = []string{}
	for i := 0; i < len(temp); i++ {
		if len(temp[i]) > 0 && temp[i][0] != '#' {
			config_t = append(config_t, temp[i])
		}
	}

	readworkitems(&config_t)
}

func readworkitems(config_t *[]string) { // Из текстового конфига в элементы
	for i := 0; i+1 < len(*config_t); i += 2 {
		var wi workitem
		wi.path = path.Clean((*config_t)[i])
		wi.operations = processoperations((*config_t)[i+1])
		config = append(config, wi)
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
