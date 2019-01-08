package coveralls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Bool(v bool) *bool        { return &v }
func Flaot(v float32) *float32 { return &v }
func URL(path, query string) string {
	u := url.URL{
		Scheme:   "https",
		Host:     "coveralls.io",
		Path:     path,
		RawQuery: query,
	}
	return u.String()
}

type GitProvider string

var (
	Github    GitProvider = "github"
	Bitbucket GitProvider = "bitbucket"
	Gitlab    GitProvider = "gitlab"
	Stash     GitProvider = "stash"
	Manual    GitProvider = "manual"
)

type Client struct {
	headers map[string]string
	client  *http.Client
}

func (c *Client) NewRequest(method, u string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %v", err)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %v", err)
	}
	if 200 > res.StatusCode || res.StatusCode >= 300 {
		return fmt.Errorf("failed status: %v %v", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return fmt.Errorf("decode: %v", err)
	}
	return nil
}

func NewClient(token string) *Client {
	c := &Client{
		headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
	}
	if token != "" {
		c.headers["Authorization"] = "token " + token
	}
	return c
}
