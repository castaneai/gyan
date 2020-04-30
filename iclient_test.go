package gyan

import (
	"os"
	"testing"
)

func newTestInternalClient(t *testing.T) *InternalClient {
	session, ok := os.LookupEnv("GYAZO_SESSION")
	if !ok {
		t.Fatal("env: GYAZO_SESSION not set")
	}
	return NewInternalClient(session)
}
