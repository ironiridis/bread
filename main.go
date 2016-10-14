package main

import "log"

import "github.com/ironiridis/tension"
import "github.com/ironiridis/private"

func main() {
	s := tension.New(private.SlackTestBotToken())
	r, err := s.AuthTest()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to Slack: %+v", r)

	// Launch websocket server, freak out if it ever returns
	go func() {
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
