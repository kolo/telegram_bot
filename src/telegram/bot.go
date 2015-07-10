package bot

import (
	"log"
	"time"
)

type HandleFunc func(*Client, *Message) error

func PollUpdates(c *Client, handler HandleFunc) error {
	offset := 0
	for {
		var err error
		updates, err := c.GetUpdates(offset, 100, 0)
		if err != nil {
			return err
		}
		if len(updates) > 0 {
			log.Printf("getUpdates: %d updates received\n", len(updates))
			for _, upd := range updates {
				if err := handler(c, upd.Message); err != nil {
					log.Printf("%v\n", err)
				}
			}
			offset = updates[len(updates)-1].ID + 1
		}

		time.Sleep(1000)
	}
}
