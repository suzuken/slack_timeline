package main

import (
	"fmt"
	"github.com/brettbuddin/victor"
	"github.com/brettbuddin/victor/pkg/chat/slackRealtime"
	"log"
	"os"
	"os/signal"
	"strings"
)

// Config is config wrapper
type Config struct {
	*victor.Config
	token string
}

// should set environment variable
// SLACK_CHANNEL_TIMELINE_ID=hogehoge
// SLACK_TOKEN=xoxp-...
func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatalln("SLACK_TOKEN is required")
	}
	tid := os.Getenv("SLACK_CHANNEL_TIMELINE_ID")
	if tid == "" {
		log.Fatalln("SLACK_CHANNEL_TIMELINE_ID is required")
	}

	ss := slackRealtime.NewConfig(token)

	config := victor.Config{
		Name:          "victor",
		ChatAdapter:   "slackRealtime",
		StoreAdapter:  "memory",
		HTTPAddr:      ":8000",
		AdapterConfig: ss,
	}
	bot := victor.New(config)

	bot.HandleCommandFunc(".*", func(s victor.State) {
		log.Printf("%v\n", s)
		if strings.HasPrefix(s.Message().ChannelName(), "_") {
			s.Chat().Send(tid, fmt.Sprintf("%s: %s", s.Message().UserName(), s.Message().Text()))
		}
	})

	go bot.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	bot.Stop()
}
