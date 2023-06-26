package main

import (
	"flag"
	"log"
)

type Conf struct {
	port uint
	root string
}

func GetCliOptions() *Conf {
	options := &Conf{}
	flag.UintVar(&options.port, "port", 8080, "the port to run the server on")
	flag.StringVar(&options.root, "root", "/", "the root folder")
	flag.Parse()

	if !isFlagPassed("port") {
		log.Println("Port unspecified, defaulting to 8080")
	}

	if !isFlagPassed("root") {
		log.Println("Root unspecified, defaulting to '/'")
	}
	return options
}
