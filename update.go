package gyan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type UpdateImageRequest struct {
	ImageID          string `json:"-"`
	Desc             string `json:"desc,omitempty"`
	MetadataIsPublic bool   `json:"metadata_is_public"`
}

func (c *InternalClient) Update(ctx context.Context, ureq *UpdateImageRequest) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	if err := enc.Encode(ureq); err != nil {
		return err
	}
	rurl := fmt.Sprintf("%s/images/%s", internalAPIEndpoint, ureq.ImageID)
	req, err := c.newRequest(ctx, "PATCH", rurl, &b)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	csrfToken, err := c.getCSRFToken(ctx, ureq.ImageID)
	if err != nil {
		return err
	}
	req.Header.Add("x-csrf-token", csrfToken)

	var resp interface{}
	if err := c.doRequest(req, &resp); err != nil {
		return err
	}
	return nil
}
