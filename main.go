package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func isFlagPassed(name string) (found bool) {
	found = false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return
}

type Conf struct {
	port uint
	root string
}

func serve(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	arg := string(reqBody)
	cmd := exec.Command("catdream", arg)
	stdout, err := cmd.Output()

	if err != nil {
		return
	}

	fmt.Fprintf(w, string(stdout))
}

func getCliOptions() (options Conf) {
	flag.UintVar(&options.port, "port", 8080, "the port to run the server on")
	flag.StringVar(&options.root, "root", "/", "the root folder")
	flag.Parse()

	if !isFlagPassed("port") {
		log.Println("Port unspecified, defaulting to 8080")
	}

	if !isFlagPassed("root") {
		log.Println("Root unspecified, defaulting to '/'")
	}
	return
}

func main() {
	options := getCliOptions()

	http.HandleFunc(options.root, serve)

	fmt.Printf("Starting server at port %d, with root '%s'\n", options.port, options.root)

	portStr := fmt.Sprintf(":%d", options.port)

	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
