package bot

import (
	"log"
	"strconv"
	"time"
)

type HandleFunc func(*Bot, *Message) error

type Bot struct {
	client *Client
}

func NewBot(token string) *Bot {
	return &Bot{
		client: NewClient(token),
	}
}

func (b *Bot) PollUpdates(handler HandleFunc) error {
	offset := 0
	for {
		var err error
		updates, err := b.GetUpdates(offset, 100, 0)
		if err != nil {
			return err
		}
		if len(updates) > 0 {
			log.Printf("getUpdates: %d updates received\n", len(updates))
			for _, upd := range updates {
				if err := handler(b, upd.Message); err != nil {
					log.Printf("%v\n", err)
				}
			}
			offset = updates[len(updates)-1].ID + 1
		}

		time.Sleep(1000)
	}
}

func (b *Bot) GetMe() (*User, error) {
	user := &User{}
	if err := b.client.Get("getMe", nil, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (b *Bot) GetUpdates(offset int, limit int, timeout int) ([]*Update, error) {
	params := RequestParams{
		"offset":  strconv.Itoa(offset),
		"limit":   strconv.Itoa(limit),
		"timeout": strconv.Itoa(timeout),
	}
	updates := []*Update{}
	if err := b.client.Get("getUpdates", params, &updates); err != nil {
		return nil, err
	}

	return updates, nil
}

func (b *Bot) SendMessage(chatID int, text string) (*Message, error) {
	params := RequestParams{
		"chat_id": strconv.Itoa(chatID),
		"text":    text,
	}
	m := &Message{}
	if err := b.client.Post("sendMessage", params, m); err != nil {
		return nil, err
	}

	return m, nil
}
