package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"telegram"
)

type Command struct {
	Name string
	Run  func(cmd *Command, args string) (string, error)
}

type Bot struct {
	Name     string
	bzClient *BugzillaClient
}

func (b *Bot) Handle(c *telegram.Client, msg *telegram.Message) error {
	if msg.Text == "" {
		return nil
	}

	if !strings.HasPrefix(msg.Text, "/") {
		return nil
	}

	args := strings.SplitN(msg.Text, " ", 2)

	cmd := args[0]
	if i := strings.Index(cmd, "@"); i != -1 {
		if cmd[i:] != b.Name {
			// Nothing to do, command was sent to other bot.
			return nil
		}
		cmd = cmd[0:i]
	}
	args = args[1:]

	var resp string
	var err error

	switch cmd {
	case "start":
		resp, err = b.start()
	case "bug":
		resp, err = b.bug(args)
	}

	if err != nil {
		return err
	}

	m, err := c.SendMessage(msg.Chat.ID, resp)
	if err != nil {
		return err
	}
	log.Printf("message (id=%d) was sent\n", m.ID)

	return nil
}

func (b *Bot) start() (string, error) {
	return renderTemplate("welcome", b.Name)
}

func (b *Bot) bug(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("command error: not enough arguments for bug command")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}

	bug, err := b.bzClient.GetBug(id)
	if err != nil {
		return "", err
	}

	return renderTemplate("bug", bug)
}
