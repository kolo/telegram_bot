package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/kolo/xmlrpc"
)

type Bug struct {
	ID       int    `xmlrpc:"id"`
	Summary  string `xmlrpc:"summary"`
	Status   string `xmlrpc:"status"`
	Assignee string `xmlrpc:"assigned_to"`
}

type BugzillaClient struct {
	client *xmlrpc.Client
	token  string
}

type params map[string]interface{}

func newBugzillaClient(url string, username string, password string) (*BugzillaClient, error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client, err := xmlrpc.NewClient(url, transport)
	if err != nil {
		return nil, err
	}

	// Get authentication token
	var resp struct {
		ID    int    `xmlrpc:"id"`
		Token string `xmlrpc:"token"`
	}
	if err = client.Call("User.login", params{"login": username, "password": password}, &resp); err != nil {
		return nil, err
	}

	return &BugzillaClient{
		client: client,
		token:  resp.Token,
	}, nil
}

func (c *BugzillaClient) GetBug(id int) (*Bug, error) {
	var resp struct {
		Bugs []Bug `xmlrpc:"bugs"`
	}
	if err := c.client.Call("Bug.get", params{"ids": []int{id}, "token": c.token}, &resp); err != nil {
		return nil, err
	}

	if len(resp.Bugs) == 0 {
		return nil, fmt.Errorf("bug %d not found", id)
	}

	return &resp.Bugs[0], nil
}
