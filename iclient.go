package gyan

import (
	"net/http"
)

type InternalClient struct {
	hc      *http.Client
	session string
}

func NewInternalClient(session string) *InternalClient {
	return &InternalClient{
		hc:      http.DefaultClient,
		session: session,
	}
}
