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
		Token    string
		Name     string
		Url      string
		Username string
		Password string
	}
)

func main() {
	flag.StringVar(&config.Token, "token", os.Getenv("bugzilla_bot_token"), "Set bot token")
	flag.StringVar(&config.Name, "password", os.Getenv("bugzilla_bot_name"), "Set bot name")
	flag.StringVar(&config.Url, "url", os.Getenv("bugzilla_bot_xmlrpc_url"), "Set Bugzilla xmlrpc endpoint url")
	flag.StringVar(&config.Username, "username", os.Getenv("bugzilla_bot_username"), "Set Bugzilla username")
	flag.StringVar(&config.Password, "password", os.Getenv("bugzilla_bot_password"), "Set Bugzilla password")

	flag.Parse()

	if config.Token == "" {
		fail(fmt.Errorf("error: token cannot be empty"))
	}

	bzClient, err := newBugzillaClient(config.Url, config.Username, config.Password)
	if err != nil {
		fail(err)
	}

	bot := &Bot{
		Name:     config.Name,
		bzClient: bzClient,
	}

	tgClient := telegram.NewClient(config.Token)
	go telegram.PollUpdates(tgClient, bot.Handle)

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
