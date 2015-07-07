package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	TelegramBaseUrl = "https://api.telegram.org"
	LimitMax        = 100
)

type requestParams map[string]string

type requestError struct {
	code int
	desc string
}

func (e requestError) Error() string {
	return fmt.Sprintf("request error: code %d, %s", e.code, e.desc)
}

type requestResult struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`

	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
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

func (c *Client) GetMe() (*User, error) {
	user := &User{}
	if err := c.get("getMe", nil, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) GetUpdates(offset int, limit int, timeout int) ([]*Update, error) {
	params := requestParams{
		"offset":  strconv.Itoa(offset),
		"limit":   strconv.Itoa(limit),
		"timeout": strconv.Itoa(timeout),
	}
	updates := []*Update{}
	if err := c.get("getUpdates", params, &updates); err != nil {
		return nil, err
	}

	return updates, nil
}

func (c *Client) SendMessage(chatID int, text string) (*Message, error) {
	params := requestParams{
		"chat_id": strconv.Itoa(chatID),
		"text":    text,
	}
	m := &Message{}
	if err := c.post("sendMessage", params, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (c *Client) get(method string, p requestParams, value interface{}) error {
	requestUrl, err := url.Parse(c.requestBaseUrl(method))
	if err != nil {
		return err
	}

	params := url.Values{}
	for key, value := range p {
		params.Add(key, value)
	}
	requestUrl.RawQuery = params.Encode()

	resp, err := c.httpClient.Get(requestUrl.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return requestError{code: resp.StatusCode, desc: "http request error"}
	}

	res := &requestResult{}
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(res); err != nil {
		return err
	}

	if !res.Ok {
		return requestError{code: res.ErrorCode, desc: res.Description}
	}

	if err = json.Unmarshal(res.Result, value); err != nil {
		return err
	}

	return nil
}

func (c *Client) post(method string, p requestParams, value interface{}) error {
	data := url.Values{}
	for key, value := range p {
		data.Add(key, value)
	}

	resp, err := c.httpClient.PostForm(c.requestBaseUrl(method), data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return requestError{code: resp.StatusCode, desc: "http request error"}
	}

	res := &requestResult{}
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(res); err != nil {
		return err
	}

	if !res.Ok {
		return requestError{code: res.ErrorCode, desc: res.Description}
	}

	if err = json.Unmarshal(res.Result, value); err != nil {
		return err
	}

	return nil
}

func (c *Client) requestBaseUrl(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", TelegramBaseUrl, c.token, method)
}
