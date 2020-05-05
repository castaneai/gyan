package gyan

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

const (
	internalAPIEndpoint = "https://gyazo.com/api/internal"
)

type SearchedImage struct {
	ID           string               `json:"image_id"`
	PermalinkURL string               `json:"permalink_url"`
	ThumbURL     string               `json:"thumb_url"`
	URL          string               `json:"url"`
	Metadata     *SearchImageMetadata `json:"metadata"`
	Owned        bool                 `json:"owned"`
	AccessPolicy string               `json:"access_policy"`
	Accessible   bool                 `json:"accessible"`
	CreatedAt    time.Time            `json:"created_at"`
	Desc         string               `json:"desc"`
}

type SearchImageMetadata struct {
	App        string `json:"app,omitempty"`
	Desc       string `json:"desc,omitempty"`
	OCRStarted bool   `json:"ocr_started,omitempty"`
	Title      string `json:"title,omitempty"`
	URL        string `json:"url,omitempty"`
}

type searchResponse struct {
	Captures []*SearchedImage `json:"captures"`
}

func (c *InternalClient) Search(ctx context.Context, query string, page int) ([]*SearchedImage, error) {
	q := url.QueryEscape(query)
	ts := float64(time.Now().UnixNano() / 1000.0)
	rurl := fmt.Sprintf(internalAPIEndpoint+"/search_result?page=%d&per=40&query=%s&timestamp=%v",
		page, q, ts)
	req, err := c.newRequest(ctx, "GET", rurl, nil)
	if err != nil {
		return nil, err
	}

	var sresp searchResponse
	if err := c.doRequest(req, &sresp); err != nil {
		return nil, err
	}
	return sresp.Captures, nil
}
