package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	TelegramBaseUrl = "https://api.telegram.org"
)

type RequestParams map[string]string

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

func (c *Client) Get(method string, p RequestParams, value interface{}) error {
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

func (c *Client) Post(method string, p RequestParams, value interface{}) error {
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
