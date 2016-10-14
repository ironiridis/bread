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
	E string `json:"event"`
	M string `json:"message,omitempty"`
}

type WSMsgClientCommand struct {
	S string // source
	C string // class
	T string // text
}

func wsserve() (err error) {
	http.HandleFunc("/breadws/mspace", wsupgrader)
	err = http.ListenAndServe(":7188", nil)
	return
}

func wsupgrader(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	rxch := make(chan *WSMsgClientEvent)
	go wsreader(conn, rxch)
	for {
		select {
		case <-time.After(time.Second * 10):
			conn.WriteJSON(
				&WSMsgClientCommand{
					S: "Server",
					C: "system",
					T: "Connection check"})
		case msg := <-rxch:
			log.Printf("received: %+v\n", msg)
			if msg == nil {
				return
			}
		}
	}
}

func wsreader(conn *websocket.Conn, rxch chan *WSMsgClientEvent) {
	defer close(rxch)
	for {
		msg := &WSMsgClientEvent{}
		err := conn.ReadJSON(msg)
		if err != nil {
			rxch <- &WSMsgClientEvent{E: "error", M: err.Error()}
			return
		}
		rxch <- msg
	}
}
