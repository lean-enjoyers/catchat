package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   int
	conn websocket.Conn
}

var clients []Client
var idToTake int = 0

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

func reader(conn *websocket.Conn, connId int) {
	clients = append(clients, Client{idToTake, *conn})
	idToTake += 1

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		for _, client := range clients {
			if client.id == connId {
				continue
			}
			if err := client.conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	// Allow all origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Client %d successfully connected...\n", idToTake)

	reader(ws, idToTake)
}

func setupRoutes(options *Conf) {
	http.HandleFunc("/say", serve)
	http.HandleFunc("/chat", serveWs)
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

	setupRoutes(&options)

	fmt.Printf("Starting server at port %d, with root '%s'\n", options.port, options.root)
	portStr := fmt.Sprintf(":%d", options.port)

	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
