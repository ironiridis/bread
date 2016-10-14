package main

import "log"

import "github.com/ironiridis/tension"
import "github.com/ironiridis/private"

var slack *tension.Slack

func main() {
	slack = tension.New(private.SlackTestBotToken())
	r, err := slack.AuthTest()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to Slack: %+v", r)

	// Launch websocket server, freak out if it ever returns
	func() {
		err := wsserve()
		panic(err)
	}()

	/*
		for {
			select {
			case rxmsg := <-rx:
				log.Printf("%s", rxmsg)
			case <-time.After(time.Second * 10):
				tx <- `{"id": 1, "type": "ping"}`
			}
		}*/
}
