package nut

import (
	"context"
)

const timeout = 1

type Client struct {
	host     string
	port     int
	username string
	password string
}

func New(host string, port int, username, password string) (*Client, error) {
	return &Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}, nil
}

func (c *Client) GetList(_ context.Context) ([]int64, error) {
	return nil, nil
}
