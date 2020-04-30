package gyan

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testImageData = "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x01\x03\x00\x00\x00%\xdbV\xca\x00\x00\x00\x06PLTE\xff\x00\x00\xff\xff\xffA\x1d4\x11\x00\x00\x00\tpHYs\x00\x00\x0e\xc4\x00\x00\x0e\xc4\x01\x95+\x0e\x1b\x00\x00\x00\nIDAT\x08\x99c`\x00\x00\x00\x02\x00\x01\xf4qd\xa6\x00\x00\x00\x00IEND\xaeB`\x82"
const testImageData2 = "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x01\x03\x00\x00\x00%\xdbV\xca\x00\x00\x00\x06PLTE\x00\xff\x00\x00\x11\xff\x11\xd8\xb2\xde\x00\x00\x00\tpHYs\x00\x00\x0e\xc4\x00\x00\x0e\xc4\x01\x95+\x0e\x1b\x00\x00\x00\nIDAT\x08\x99c`\x00\x00\x00\x02\x00\x01\xf4qd\xa6\x00\x00\x00\x00IEND\xaeB`\x82"

func TestUpload(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	testCases := []struct {
		name    string
		request *UploadRequest
	}{
		{name: "SimpleFileUpload", request: &UploadRequest{
			Filename:  fmt.Sprintf("SimpleFileUpload%v", time.Now().UnixNano()),
			ImageData: bytes.NewReader([]byte(testImageData)),
		}},
		{name: "UploadWithMetadata", request: &UploadRequest{
			Filename:         fmt.Sprintf("UploadWithMetadata%v", time.Now().UnixNano()),
			ImageData:        bytes.NewReader([]byte(testImageData2)),
			MetadataIsPublic: true,
			RefererURL:       "http://gyan-test.example",
			Title:            "gyan-test-image-title",
			Desc:             "gyan-test-image-desc #tag1 #tag2",
			CreatedAt:        time.Now().Add(-1 * time.Hour),
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img, err := c.Upload(ctx, tc.request)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("%+v", img)
			assert.NotEmpty(t, img.ID)
			assert.NotEmpty(t, img.URL)
		})
	}

}
