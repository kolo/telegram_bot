package main

import "fmt"

type GroupChat struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Message struct {
	ID   int   `json:"message_id"`
	From *User `json:"from"`
	Date int   `json:"date"`

	Chat struct {
		ID int `json:"id"`
	} `json:"chat"`

	Text string `json:"text"`
}

func (m *Message) String() string {
	return fmt.Sprintf("{id: %d, from:%v, date: %d, chat: %d, text: %s",
		m.ID, *m.From, m.Date, m.Chat.ID, m.Text)
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message"`
}
