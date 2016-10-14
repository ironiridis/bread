package main

// This is going to be a very poorly factored golang app. MVP and all that, you know.

import "log"
import "time"
import "github.com/ironiridis/tension"
import "github.com/ironiridis/private"

func main() {
	err := wsserve()
	if err != nil {
		panic(err)
	}

	if true {
		return
	}
	s := tension.New(private.SlackTestToken())
	tx, rx, err := slackRTMWebsocket(s)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected, presumably")

	for {
		select {
		case rxmsg := <-rx:
			log.Printf("%s", rxmsg)
		case <-time.After(time.Second * 10):
			tx <- `{"id": 1, "type": "ping"}`
		}
	}
}
