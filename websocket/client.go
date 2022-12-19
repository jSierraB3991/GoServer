package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub      *Hub
	id       string
	socket   *websocket.Conn
	outbound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (this *Client) Write() {
	for {
		select {
		case message, ok := <-this.outbound:
			if !ok {
				this.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			this.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
