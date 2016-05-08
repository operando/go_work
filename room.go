package main

import (
	"chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte      // forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	join    chan *client     // joinはチャットルームに参加しようとしているクライアントのためのチャネル
	leave   chan *client     // leaveはチャットルームから退室しようとしているクライアントのためのチャネル
	clients map[*client]bool // clientsには在室しているすべてのクライアントが保持される
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("join new clint")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("leave clint")
		case msg := <-r.forward:
			r.tracer.Trace("send message : ", string(msg))
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
					r.tracer.Trace("send done")
				default:
					// fail send message
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("fail sned")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

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
	go client.write()
	client.read()
}
