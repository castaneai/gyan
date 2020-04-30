package gyan

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	c := newTestInternalClient(t)
	ctx := context.Background()
	images, err := c.Search(ctx, "test", 1)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", images)
	assert.NotNil(t, images)
}
