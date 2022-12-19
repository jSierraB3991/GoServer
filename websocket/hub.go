package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client, 0),
		unregister: make(chan *Client, 0),
		mutex:      &sync.Mutex{},
	}
}

func (this *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	client := NewClient(this, socket)
	this.register <- client

	go client.Write()
}

func (this *Hub) Run() {
	for {
		select {
		case client := <-this.register:
			this.OnConnect(client)
		case client := <-this.unregister:
			this.OnDisconnect(client)

		}
	}
}

func (this *Hub) OnConnect(client *Client) {
	log.Println("A new Client Connected", client.socket.RemoteAddr())
	this.mutex.Lock()
	defer this.mutex.Unlock()

	client.id = client.socket.RemoteAddr().String()
	this.clients = append(this.clients, client)
}

func (this *Hub) OnDisconnect(client *Client) {
	log.Println("A new Client DiscConnected", client.socket.RemoteAddr())
	client.socket.Close()
	this.mutex.Lock()
	defer this.mutex.Unlock()

	i := -1
	for index, c := range this.clients {
		if c.id == client.id {
			i = index
		}
	}

	copy(this.clients[i:], this.clients[i+1:])
	this.clients[len(this.clients)-1] = nil
	this.clients = this.clients[:len(this.clients)-1]
}

func (this *Hub) Broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range this.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}
