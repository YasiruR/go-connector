package http

import (
	"bytes"
	"fmt"
	"github.com/YasiruR/connector/core/errors"
	"io"
	"net/http"
)

type Client struct {
	hc *http.Client
}

func NewClient() *Client {
	return &Client{hc: http.DefaultClient}
}

func (c *Client) Send(data []byte, destination string) (response []byte, err error) {
	if data == nil {
		return nil, fmt.Errorf("GET method is not implemented yet")
	}

	response, status, err := c.post(destination, data)
	if err != nil {
		return nil, errors.SendFailed(destination, http.MethodPost, err)
	}

	if status != http.StatusOK && status != http.StatusCreated {
		return nil, errors.InvalidStatusCode(status)
	}

	return response, nil
}

func (c *Client) post(url string, data []byte) (resData []byte, statusCode int, err error) {
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
