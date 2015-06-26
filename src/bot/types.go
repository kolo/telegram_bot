package main

import "encoding/json"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Update struct {
	ID      int             `json:"update_id"`
	Message json.RawMessage `json:"message"`
}
