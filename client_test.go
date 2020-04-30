package gyan

import (
	"os"
	"testing"
)

func newTestClient(t *testing.T) *Client {
	token, ok := os.LookupEnv("GYAZO_TOKEN")
	if !ok {
		t.Fatal("env: GYAZO_TOKEN not set")
	}
	return NewClient(token)
}
