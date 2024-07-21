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

func (c *Client) Post(url string, data []byte) (resData []byte, statusCode int, err error) {
	res, err := c.hc.Post(url, `application/json`, bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer res.Body.Close()

	resData, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body: %w", err)
	}

	return resData, res.StatusCode, nil
}
