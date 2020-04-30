package gyan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type UploadRequest struct {
	Filename         string
	ImageData        io.Reader `json:"imagedata"`
	MetadataIsPublic bool      `json:"metadata_is_public"`
	RefererURL       string    `json:"referer_url"`
	Title            string    `json:"title"`
	Desc             string    `json:"desc"`
	CreatedAt        time.Time `json:"created_at"`
	CollectionID     string    `json:"collection_id"`
}

type UploadedImage struct {
	ID           string `json:"image_id"`
	PermalinkURL string `json:"permalink_url"`
	ThumbURL     string `json:"thumb_url"`
	URL          string `json:"url"`
	Type         string `json:"type"`
}

func (c *Client) Upload(ctx context.Context, request *UploadRequest) (*UploadedImage, error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile("imagedata", request.Filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fw, request.ImageData); err != nil {
		return nil, err
	}

	if request.MetadataIsPublic {
		if err := writeMultipartFormField(mw, "metadata_is_public", "true"); err != nil {
			return nil, err
		}
	}
	if request.RefererURL != "" {
		if err := writeMultipartFormField(mw, "referer_url", request.RefererURL); err != nil {
			return nil, err
		}
	}
	if request.Title != "" {
		if err := writeMultipartFormField(mw, "title", request.Title); err != nil {
			return nil, err
		}
	}
	if request.Desc != "" {
		if err := writeMultipartFormField(mw, "desc", request.Desc); err != nil {
			return nil, err
		}
	}
	if !request.CreatedAt.IsZero() {
		if err := writeMultipartFormField(mw, "created_at", fmt.Sprintf("%f", float64(request.CreatedAt.Unix()))); err != nil {
			return nil, err
		}
	}
	if request.CollectionID != "" {
		if err := writeMultipartFormField(mw, "collection_id", request.CollectionID); err != nil {
			return nil, err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, err
	}

	url := uploadEndpoint + "/api/upload"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %w", err)
	}
	req.Header.Add("Content-Type", mw.FormDataContentType())
	req = req.WithContext(ctx)

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, newResponseError(res)
	}

	img := &UploadedImage{}
	if err := json.NewDecoder(res.Body).Decode(img); err != nil {
		return nil, err
	}
	return img, nil
}

func writeMultipartFormField(w *multipart.Writer, key, value string) error {
	field, err := w.CreateFormField(key)
	if err != nil {
		return err
	}
	if _, err := field.Write([]byte(value)); err != nil {
		return err
	}
	return nil
}
