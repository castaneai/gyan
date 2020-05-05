package gyan

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

var regexCSRFToken = regexp.MustCompile(`<meta name="csrf-token" content="([0-9A-Za-z+/=]+)" />`)

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

func (c *InternalClient) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
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
	return req, nil
}

func (c *InternalClient) doRequest(req *http.Request, ress interface{}) error {
	res, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		errb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			errb = []byte("<unknown>")
		}
		return fmt.Errorf("http error: %v: %+v", res.Status, errb)
	}
	if err := json.NewDecoder(res.Body).Decode(&ress); err != nil {
		return err
	}
	return nil
}

func (c *InternalClient) getCSRFToken(ctx context.Context, imageID string) (string, error) {
	rurl := "https://gyazo.com/" + imageID
	req, err := c.newRequest(ctx, "GET", rurl, nil)
	if err != nil {
		return "", err
	}
	res, err := c.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		errb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			errb = []byte("<unknown>")
		}
		return "", fmt.Errorf("http error: %v: %+v", res.Status, errb)
	}

	rbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	ms := regexCSRFToken.FindStringSubmatch(string(rbody))
	if len(ms) < 2 {
		return "", fmt.Errorf("response body does not match pattern: %v", regexCSRFToken)
	}
	return ms[1], nil
}
