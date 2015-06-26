package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	token string
)

func init() {
}

func main() {
	flag.StringVar(&token, "token", os.Getenv("bugzilla_bot_token"), "Set bot token")
	flag.Parse()

	if token == "" {
		fail(fmt.Errorf("error: token cannot be empty"))
	}

	client := NewClient(token)
	user, err := client.GetMe()
	if err != nil {
		fail(err)
	}

	log.Printf("%s starts working\n", user.Username)

	go pollCommands(client)

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
	log.Fatalf("%v\n", err)
}

func pollCommands(client *Client) {
	offset := 0
	for {
		var err error
		updates, err := client.GetUpdates(offset, 100, 0)
		if err != nil {
			fail(err)
		}
		if len(updates) > 0 {
			log.Printf("getUpdates: %d updates received\n", len(updates))
			for _, upd := range updates {
				log.Printf("update %d: %s\n", upd.ID, upd.Message)
			}
			offset = updates[len(updates)-1].ID + 1
		}

		time.Sleep(1000)
	}
}
