package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	TelegramBaseUrl = "https://api.telegram.org"
	LimitMax        = 100
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	httpClient := &http.Client{}
	return &Client{
		token:      token,
		httpClient: httpClient,
	}
}

func (c *Client) GetMe() error {
	baseUrl, err := url.Parse(c.requestBaseUrl("getMe"))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Get(baseUrl.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	return nil
}

func (c *Client) requestBaseUrl(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", TelegramBaseUrl, token, method)
}
