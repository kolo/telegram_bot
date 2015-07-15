package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"telegram"
)

var (
	config struct {
		Token string
	}
)

func main() {
	flag.StringVar(&config.Token, "token", os.Getenv("news_bot_token"), "Set bot token")
	flag.Parse()

	if config.Token == "" {
		fail(fmt.Errorf("error: token cannot be empty"))
	}

	c := telegram.NewClient(config.Token)
	go telegram.PollUpdates(c, handle)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	cleanup := make(chan bool)
	go func() {
		for sig := range signalChan {
			log.Printf("signal: %v\n", sig)
			cleanup <- true
		}
	}()
	<-cleanup
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func handle(c *telegram.Client, msg *telegram.Message) error {
	m, err := c.SendMessage(msg.Chat.ID, msg.Text)
	if err != nil {
		return err
	}
	log.Printf("message (id:%d) sent\n", m.ID)

	return nil
}
