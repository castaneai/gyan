package gyan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	internalAPIEndpoint = "https://gyazo.com/api/internal"
)

type SearchedImage struct {
	ID           string         `json:"image_id"`
	PermalinkURL string         `json:"permalink_url"`
	ThumbURL     string         `json:"thumb_url"`
	URL          string         `json:"url"`
	Metadata     *ImageMetadata `json:"metadata"`
	Owned        bool           `json:"owned"`
	AccessPolicy string         `json:"access_policy"`
	Accessible   bool           `json:"accessible"`
	CreatedAt    time.Time      `json:"created_at"`
	Desc         string
}

type ImageMetadata struct {
	App        string `json:"app"`
	Desc       string `json:"desc"`
	OCRStarted bool   `json:"ocr_started"`
	Title      string `json:"title"`
	URL        string `json:"url"`
}

type searchResponse struct {
	Captures []*SearchedImage `json:"captures"`
}

func (c *InternalClient) Search(ctx context.Context, query string, page int) ([]*SearchedImage, error) {
	q := url.QueryEscape(query)
	ts := float64(time.Now().UnixNano() / 1000.0)
	rurl := fmt.Sprintf(internalAPIEndpoint+"/search_result?page=%d&per=40&query=%s&timestamp=%v",
		page, q, ts)
	req, err := http.NewRequest("GET", rurl, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.AddCookie(&http.Cookie{
		Name:     "Gyazo_session",
		Value:    c.session,
		HttpOnly: true,
		Secure:   true,
	})

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		errb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			errb = []byte("<unknown>")
		}
		return nil, fmt.Errorf("http error: %v: %+v", res.Status, errb)
	}

	var sresp searchResponse
	if err := json.NewDecoder(res.Body).Decode(&sresp); err != nil {
		return nil, err
	}
	return sresp.Captures, nil
}
