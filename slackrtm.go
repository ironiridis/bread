package main

import "log"
import "bytes"
import "golang.org/x/net/websocket"
import "github.com/ironiridis/tension"

func slackRTMWorker(ws *websocket.Conn, tx, rx chan string) {
	go func() {
		rbuf := make([]byte, 4000)
		for {
			n, err := ws.Read(rbuf)
			if err != nil {
				close(rx)
				log.Fatalf("slack websocket rx error: %#v", err)
				return
			}
			rx <- string(rbuf[:n])
		}
	}()
	go func() {
		// we use bytes.Buffer to utilize the automatic partial write
		// retry logic it has baked in.
		tbuf := bytes.Buffer{}
		for txstr := range tx {
			tbuf.WriteString(txstr)
			_, err := tbuf.WriteTo(ws)
			if err != nil {
				log.Fatalf("slack websocket tx error: %#v", err)
				return
			}
			tbuf.Reset()
		}
	}()
}

func slackRTMWebsocket(s *tension.Slack) (tx, rx chan string, err error) {
	rtm, err := s.RTMStart(true, true)
	if err != nil {
		return
	}

	origin := "http://localhost/"
	ws, err := websocket.Dial(rtm.URL, "", origin)
	if err != nil {
		return
	}

	tx = make(chan string, 5)
	rx = make(chan string, 5)

	go slackRTMWorker(ws, tx, rx)
	return
}
