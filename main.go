package main

import "log"

import "github.com/ironiridis/tension"
import "github.com/ironiridis/private"

var slackAPI *tension.Slack
var slackBot *tension.Slack

func checkerr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	slackAPI = tension.New(private.SlackTestAccessToken())
	r, err := slackAPI.AuthTest()
	checkerr(err)
	log.Printf("api authed to Slack: %+v", r)

	slackBot = tension.New(private.SlackTestBotToken())
	r, err = slackBot.AuthTest()
	checkerr(err)
	log.Printf("bot authed to Slack: %+v", r)

	rtmresult, err := slackBot.RTMStart(true, true)
	checkerr(err)
	log.Printf("bot rtm setup with Slack: %+v", rtmresult)

	rtm, err := rtmresult.Dial()
	checkerr(err)
	log.Printf("bot rtm connected with Slack: %+v", rtm)

	// Launch websocket server, freak out if it ever returns
	go func() {
		err := wsserve()
		panic(err)
	}()

	for m := range rtm.Rx {
		log.Println(m)
	}

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
