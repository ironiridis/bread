package main

import (
	"log"
	"math/rand"
)

type Msg struct {
	Source string
	Text   string
}

const SlackChannelNameLen = 21
const SlackChannelNameBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func slackChanName(prefix string) string {
	r := make([]byte, 0, SlackChannelNameLen)
	r = append(r, prefix...)
	r = append(r, "_"...)

	for len(r) < SlackChannelNameLen {
		r = append(r, SlackChannelNameBytes[rand.Intn(len(SlackChannelNameBytes))])
	}
	return string(r)
}

func slackGetChannel() (tx, rx chan *Msg, err error) {
	r, err := slackAPI.GroupCreate(slackChanName("bread"))
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
