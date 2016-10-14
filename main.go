package main

import "log"

import "github.com/ironiridis/tension"
import "github.com/ironiridis/private"

var slackAPI *tension.Slack
var slackBot *tension.Slack

func main() {
	slackBot = tension.New(private.SlackTestBotToken())
	r, err := slackBot.AuthTest()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bot connected to Slack: %+v", r)
	slackAPI = tension.New(private.SlackTestToken())
	r, err = slackAPI.AuthTest()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("api connected to Slack: %+v", r)

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
