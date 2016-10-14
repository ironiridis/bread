package main

import (
	"fmt"
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

	toSlack, fromSlack, err := slackGetChannel()
	if err != nil {
		log.Printf("Slack returned an error during a WS setup: %v", err)
		return
	}

	rxch := make(chan *WSMsgClientEvent, 2)
	go wsreader(conn, rxch)
	for {
		select {
		case <-time.After(time.Second * 10):
			err = conn.WriteJSON(
				&WSMsgClientCommand{
					S: "Server",
					C: "system",
					T: "timeout"})
		case msg := <-fromSlack:
			if msg == nil {
				log.Println("Slack is gone with an active WS connection")
				return
			}
			err = conn.WriteJSON(&WSMsgClientCommand{
				S: msg.Source,
				C: "received",
				T: msg.Text})
		case msg := <-rxch:
			if msg == nil {
				return
			}
			switch msg.E {
			case "hello":
				err = conn.WriteJSON(&WSMsgClientCommand{
					S: "Server",
					C: "system",
					T: "hello"})
			case "ping":
				err = conn.WriteJSON(&WSMsgClientCommand{
					S: "Server",
					C: "system",
					T: "pong"})
			case "error":
				err = fmt.Errorf("%s", msg.M)
			case "typing":
				log.Println("not handling typing event: unimplemented")
			case "message":
				toSlack <- &Msg{Source: "Client", Text: msg.M}
			default:
				log.Printf("unhandled: %v\n", msg)
			}
		}
		if err != nil {
			log.Printf("wsupgrader error %v", err)
			return
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
