package main

import (
	"flag"
	"fmt"
	"os"
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
	if err := client.GetMe(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
