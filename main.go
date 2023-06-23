package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"text/template"

	"github.com/gorilla/websocket"
)

type Client struct {
	id    int
	conn  websocket.Conn
	valid bool
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

var options Conf

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

var t *template.Template
var indexTemplate = template.Must(template.ParseFiles("tmpl/index.html"))

func serveIndex(w http.ResponseWriter, r *http.Request) {
	err := t.Execute(w, map[string]interface{}{
		"port": options.port,
	})

	if err != nil {
		log.Print(err)
	}
}

func reader(conn *websocket.Conn, connId int) {

	clientContext := Client{
		id:    idToTake,
		conn:  *conn,
		valid: false,
	}
	clients = append(clients, clientContext)
	idToTake += 1

	if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("Enter your name: "))); err != nil {
		log.Println(err)
		return
	}

	_, name, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	nameStr := string(name)
	prefix := []byte(nameStr + ": ")

	// For user to be able to tell that they have joined
	if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("============="))); err != nil {
		log.Println(err)
		return
	}

	// Alert all clients that you have connected
	clients[connId].valid = true
	for _, client := range clients {
		if client.id == connId || !client.valid {
			continue
		}
		if err := client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s connected", nameStr))); err != nil {
			client.valid = false
		}
	}

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			for _, client := range clients {
				if client.id == connId || !client.valid {
					continue
				}
				if err := client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s disconnected", nameStr))); err != nil {
					client.valid = false
				}
			}
			return
		}

		for _, client := range clients {
			if !client.valid {
				continue
			}
			if err := client.conn.WriteMessage(messageType, append(prefix, p...)); err != nil {
				client.valid = false
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
	http.HandleFunc("/", serveIndex)
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
	options = getCliOptions()
	t = template.Must(template.Must(indexTemplate.Clone()).ParseFiles("tmpl/i.html"))

	setupRoutes(&options)

	fmt.Printf("Starting server at port %d, with root '%s'\n", options.port, options.root)
	portStr := fmt.Sprintf(":%d", options.port)

	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
