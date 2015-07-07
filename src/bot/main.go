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
	config struct {
		Token    string
		Url      string
		Username string
		Password string
	}

	bzClient *BugzillaClient
)

func init() {
}

func main() {
	flag.StringVar(&config.Token, "token", os.Getenv("bugzilla_bot_token"), "Set bot token")
	flag.StringVar(&config.Url, "url", os.Getenv("bugzilla_xmlrpc_url"), "Set Bugzilla xmlrpc endpoint url")
	flag.StringVar(&config.Username, "username", os.Getenv("bugzilla_username"), "Set Bugzilla username")
	flag.StringVar(&config.Password, "password", os.Getenv("bugzilla_password"), "Set Bugzilla password")

	flag.Parse()

	if config.Token == "" {
		fail(fmt.Errorf("error: token cannot be empty"))
	}

	var err error
	bzClient, err = newBugzillaClient(config.Url, config.Username, config.Password)
	if err != nil {
		fail(err)
	}

	tgClient := NewClient(config.Token)
	user, err := tgClient.GetMe()
	if err != nil {
		fail(err)
	}

	log.Printf("%s starts working\n", user.Username)

	go pollUpdates(tgClient)

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

func pollUpdates(client *Client) {
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
				if err := handle(client, upd.Message); err != nil {
					log.Printf("%v\n", err)
				}
			}
			offset = updates[len(updates)-1].ID + 1
		}

		time.Sleep(1000)
	}
}
