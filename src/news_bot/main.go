package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bot"
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

	b := bot.NewBot(config.Token)
	go b.PollUpdates(handle)

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

func handle(b *bot.Bot, msg *bot.Message) error {
	m, err := b.SendMessage(msg.Chat.ID, msg.Text)
	if err != nil {
		return err
	}
	log.Printf("message (id:%d) sent\n", m.ID)

	return nil
}
