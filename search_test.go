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
	for _, img := range images {
		log.Printf("%+v", img)
	}
	assert.NotNil(t, images)
}
