package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"utils/utils"
)

type Server struct {
	cache *utils.LRUCache
	clients map[*websocket.Conn]bool // Connected clients
	broadcast chan Message 
}

type Message struct {
	Content   string `json:"content,omitempty"`
}

type ResponseBody struct {
	Data string `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
			return true
	},
}

func NewServer() *Server {
	return &Server{
		cache: utils.CreateLRUCache(),
		clients: make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/hello" {
// 			http.Error(w, "404 not found.", http.StatusNotFound)
// 			return
// 	}

// 	if r.Method != "GET" {
// 			http.Error(w, "Method is not supported.", http.StatusNotFound)
// 			return
// 	}


// 	fmt.Fprintf(w, "Hello!")
// }

func (server *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
			log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	server.clients[ws] = true
	str, err := json.Marshal(server.cache.ToArray())
	if err != nil {
		return
	}
	server.broadcast <- Message{Content: string(str)}

	for {
			var msg Message
			// Read in a new message as JSON and map it to a Message object
			err := ws.ReadJSON(&msg)
			if err != nil {
					log.Printf("error: %v", err)
					delete(server.clients, ws)
					break
			}
			// Send the newly received message to the broadcast channel
			server.broadcast <- msg
	}
}

func (server *Server) handleMessages() {
	for {
			// Grab the next message from the broadcast channel
			msg := <-server.broadcast
			// Send it out to every client that is currently connected
			for client := range server.clients {
					err := client.WriteJSON(msg)
					if err != nil {
							log.Printf("error: %v", err)
							client.Close()
							delete(server.clients, client)
					}
			}
	}
}


func (server *Server) cacheHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	if r.URL.Path != "/cache" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "OPTIONS" { // for pre-flight requests
		return
	} else if r.Method == "GET" {
			key := r.URL.Query().Get("key")
			// fmt.Printf("key: %s", key)
			if value, ok := server.cache.Get(key); ok {
				str, err := json.Marshal(server.cache.ToArray())
				if err != nil {
					return
				}
				server.broadcast <- Message{Content: string(str)}

				res := ResponseBody{
					Data: value,
				}
				jsonResponse, err := json.Marshal(res)
				if err != nil {
						http.Error(w, "Internal server error", http.StatusInternalServerError)
						return
				}

				w.Header().Set("Content-Type", "application/json")

				// Write the JSON response to the client
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResponse)
				return;
			}
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusNotFound)
	} else if r.Method == "PUT" {
		var requestBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
		}

		key := requestBody["key"]
		value := requestBody["value"]
		duration := requestBody["duration"]
		// fmt.Printf("values: %s", requestBody["key"])
		expSec, err := strconv.Atoi(duration)
		if err != nil {
			http.Error(w, "500 server error!", http.StatusNotFound)
			return
		}
		server.cache.Put(key, value, expSec, func(msg string) {server.broadcast <- Message{Content: msg}});
		str, err := json.Marshal(server.cache.ToArray())
		if err != nil {
			w.WriteHeader(500)
			return
		}
		server.broadcast <- Message{Content: string(str)}
		w.WriteHeader(http.StatusOK)
			// fmt.Printf("cache: %s", string(str))
	} else if r.Method == "DELETE" {
		key := r.URL.Query().Get("key")
		// fmt.Printf("key: %s", key)
		if ok := server.cache.Delete(key); ok {
			// fmt.Fprintf(w, "key: %s deleted", key)
			str, err := json.Marshal(server.cache.ToArray())
			if err != nil {
				w.WriteHeader(500)
				return
			}
			server.broadcast <- Message{Content: string(str)}
			return;
		}
		w.WriteHeader(http.StatusOK)
		// fmt.Fprintf(w, "not found!")
	} else {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
	}
}

func main() {
	// init cache
	server := NewServer()

	var mu sync.Mutex

	go func() {
		for {
			mu.Lock()
			// evict
			server.cache.CheckForExp(func (str string) {
				server.broadcast <- Message{Content: string(str)}
			})
			mu.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()

	
	http.HandleFunc("/ws", server.handleConnections)
	http.HandleFunc("/cache", server.cacheHandler)
	
	go server.handleMessages()

	fmt.Printf("Starting server at port 8080\n");
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}