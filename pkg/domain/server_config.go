package domain

import (
	"flag"
	"log"

	"github.com/lean-enjoyers/catchat/pkg/cli"
)

type Conf struct {
	Port uint
	Root string
}

func GetCliOptions() *Conf {
	options := &Conf{}
	flag.UintVar(&options.Port, "port", 8080, "the port to run the server on")
	flag.StringVar(&options.Root, "root", "/", "the root folder")
	flag.Parse()

	if !cli.IsFlagPassed("port") {
		log.Println("Port unspecified, defaulting to 8080")
	}

	if !cli.IsFlagPassed("root") {
		log.Println("Root unspecified, defaulting to '/'")
	}
	return options
}
