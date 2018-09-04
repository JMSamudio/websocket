package main

import (
	"log"
	"net/http"
	"websocket/trace"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte

	// join is a channel for clients wishing to join the room.
	join chan *client

	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	// clients holds all current clients in this room.
	clients map[*client]bool

	// tracer will receive trace information of activity
	// in the room.
	tracer trace.Tracer
}

//newRoom makes a new room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//add
			r.clients[client] = true
			r.tracer.Trace("New client join")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Mess rec: ", string(msg))
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("-- send to client")

			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

//upgrade the http connection
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	//other goroutine
	go client.write()
	client.read()
}
