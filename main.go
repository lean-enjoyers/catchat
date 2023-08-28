package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"text/template"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"

	_ "github.com/lean-enjoyers/catchat/pkg/command/commands"
	"github.com/lean-enjoyers/catchat/pkg/domain"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var t *template.Template
var indexTemplate = template.Must(template.ParseFiles("tmpl/index.html"))
var store = sessions.NewCookieStore([]byte("super_secret_key"))
var users map[string]string = make(map[string]string)

func setupRoutes(hub *domain.Hub, options *domain.Conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveIndex(options, w, r)
	})
	http.HandleFunc("/say", serve)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/sign-up", registerHandler)
}

func setupTemplate() {
	t = template.Must(template.Must(indexTemplate.Clone()).ParseFiles("tmpl/chat.html"))
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

func serveIndex(options *domain.Conf, w http.ResponseWriter, r *http.Request) {
	err := t.Execute(w, map[string]interface{}{
		"port": options.Port,
	})

	if err != nil {
		log.Print(err)
	}
}

// Serve web sockets
func serveWs(hub *domain.Hub, w http.ResponseWriter, r *http.Request) {
	// Allow all origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := domain.MakeClient(conn)

	// Select the hub to connect to.
	client.SelectHub(hub)

	// Connect the client to the hub.
	client.Connect()

	go client.SendLoop()
	go client.ReceiveLoop()
}

///////////////////////////////
// Note(Appy): Authentication

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}

	// Retrieve username and password from the form
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	pswd, exists := users[username]
	if exists {
		// Returns a new session if there is no current session.
		session, _ := store.Get(r, "session.id")
		if pswd == password {
			session.Values["authenticated"] = true
			session.Save(r, w)
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		w.Write([]byte("Login successfully!"))
	} else {
		http.Error(w, "User not found", http.StatusUnauthorized)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}

	// Retrieve username and password from the form
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	confirmedPassword := r.Form.Get("confirmed_password")
	_, exists := users[username]
	if exists {
		// Returns a new session if there is no current session.
		http.Error(w, "User already exists", http.StatusUnauthorized)
	} else {
		session, _ := store.Get(r, "session.id")
		if len(username) > 0 && len(password) > 0 && len(confirmedPassword) > 0 {
			println(username)
			println(password)
			println(confirmedPassword)
			if confirmedPassword == password {
				users[username] = password
				session.Values["authenticated"] = true
				session.Save(r, w)
			} else {
				http.Error(w, "Passwords don't match", http.StatusUnauthorized)
				return
			}
			w.Write([]byte("User has been registered successfully!"))
		} else {

			http.Error(w, "Please fill in required information", http.StatusUnauthorized)
			return
		}

	}
}

func main() {
	options := domain.GetCliOptions()
	setupTemplate()

	hub := domain.MakeHub()
	go hub.Run()

	setupRoutes(hub, options)
	fmt.Printf("Starting server at port %d, with root '%s'\n", options.Port, options.Root)
	portStr := fmt.Sprintf(":%d", options.Port)

	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
