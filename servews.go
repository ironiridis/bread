package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSMsgClientEvent struct {
	E string
	M string `json:",omitempty"`
}

type WSMsgClientCommand struct {
	S string // source
	C string // class
	T string // text
}

func wsupgrader(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	rxch := make(chan *WSMsgClientEvent)
	go wsreader(conn, rxch)
	for {
		select {
		case <-time.After(time.Second * 5):
			conn.WriteJSON(&WSMsgClientCommand{S: "Pinger", C: "event", T: "ping"})
		case msg := <-rxch:
			log.Printf("received: %+v\n", msg)
		}
	}
}

func wsreader(conn *websocket.Conn, rxch chan *WSMsgClientEvent) {
	for {
		msg := &WSMsgClientEvent{}
		err := conn.ReadJSON(msg)
		if err != nil {
			msg.E = "wserror"
			msg.M = err.Error()
			rxch <- msg
			close(rxch)
		}
	}
}
