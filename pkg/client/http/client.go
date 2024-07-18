package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	hc *http.Client
}

func NewClient() *Client {
	return &Client{hc: http.DefaultClient}
}

func (c *Client) Post(url string, data []byte) (statusCode int, resData []byte, err error) {
	res, err := c.hc.Post(url, `application/json`, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return res.StatusCode, body, nil
}
