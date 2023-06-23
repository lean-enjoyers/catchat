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

var clients Hub
var idToTake int = 0
var options Conf
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var t *template.Template
var indexTemplate = template.Must(template.ParseFiles("tmpl/index.html"))

func isFlagPassed(name string) (found bool) {
	found = false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return
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
	clients.addClient(clientContext)

	idToTake += 1

	if err := clients.writeStringToClient(connId, "Enter your name: "); err != nil {
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
	if err := clients.writeStringToClient(connId, "============="); err != nil {
		log.Println(err)
		return
	}

	// Alert all clients that you have connected
	clients.setValid(connId, true)
	for _, client := range clients.clients {
		if client.id == connId || !client.valid {
			continue
		}
		if err := client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Server: %s connected", nameStr))); err != nil {
			client.valid = false
		}
	}

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			for cidx, client := range clients.clients {
				if client.id == connId || !client.valid {
					continue
				}
				if err := clients.writeStringToClient(cidx, fmt.Sprintf("%s disconnected", nameStr)); err != nil {
					client.valid = false
				}
			}
			return
		}

		for _, client := range clients.clients {
			if !client.valid {
				continue
			}
			if err := client.conn.WriteMessage(messageType, append(prefix, p...)); err != nil {
				client.valid = false
			}
		}
	}
}

// Serve web sockets
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

func main() {
	options = GetCliOptions()
	t = template.Must(template.Must(indexTemplate.Clone()).ParseFiles("tmpl/chat.html"))

	setupRoutes(&options)

	fmt.Printf("Starting server at port %d, with root '%s'\n", options.port, options.root)
	portStr := fmt.Sprintf(":%d", options.port)

	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
