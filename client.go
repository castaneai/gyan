package gyan

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	apiEndpoint = "https://api.gyazo.com"
	uploadEndpoint = "https://upload.gyazo.com"
)

type Client struct {
	hc *http.Client
}

func NewClient(token string) *Client {
	oauthClient := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	)
	return &Client{
		hc: oauthClient,
	}
}

func (c *Client) SetHttpClient(hc *http.Client) {
	c.hc = hc
}
