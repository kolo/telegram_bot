package main

import (
	"log"
	"regexp"
	"strings"
)

type Command struct {
	Name string
	Run  func(cmd *Command, args string) (string, error)
}

var (
	commandRx = regexp.MustCompile(`^\/\w+`)
	commands  = []*Command{
		cmdBug,
	}
)

func handle(client *Client, msg *Message) error {
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

			m, err := client.SendMessage(msg.Chat.ID, resp)
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
