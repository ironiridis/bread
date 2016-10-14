package main

import "log"

type Msg struct {
	Source string
	Text   string
}

func slackGetChannel() (tx, rx chan *Msg, err error) {
	r, err := slack.GroupCreate("bread")
	log.Printf("get group create result: %+v", r)
	if err != nil {
		return
	}

	tx = make(chan *Msg)
	rx = make(chan *Msg)

	go func() {
		for t := range tx {
			log.Printf("Would send %+v to Slack...", t)
		}
	}()

	return
}
