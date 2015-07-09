package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"bot"
)

type Command struct {
	Name string
	Run  func(cmd *Command, args string) (string, error)
}

var (
	config struct {
		Token    string
		Url      string
		Username string
		Password string
	}

	commandRx = regexp.MustCompile(`^\/\w+`)
	commands  = []*Command{
		cmdBug,
	}

	bzClient *BugzillaClient
)

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
	log.Fatalf("%v\n", err)
}

func handle(b *bot.Bot, msg *bot.Message) error {
	if msg.Text == "" {
		return nil
	}

	if !commandRx.MatchString(msg.Text) {
		return nil
	}

	name, args := parseCommand(msg.Text)
	for _, cmd := range commands {
		if cmd.Name == name {
			resp, err := cmd.Run(cmd, args)
			if err != nil {
				return err
			}

			m, err := b.SendMessage(msg.Chat.ID, resp)
			if err != nil {
				return err
			}
			log.Printf("message %d sent successfully\n", m.ID)
		}
	}

	return nil
}

func parseCommand(s string) (string, string) {
	var cmd, args string

	t := strings.SplitN(s, " ", 2)
	cmd = strings.TrimLeft(t[0], "/")
	if len(t) == 2 {
		args = t[1]
	}

	return cmd, args
}
