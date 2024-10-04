package http

import (
	"bytes"
	"fmt"
	"github.com/YasiruR/go-connector/domain/pkg"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	hc *http.Client
}

func NewClient(log pkg.Log) *Client {
	//rootCert, err := os.ReadFile(`certs/rootCA.crt`)
	//if err != nil {
	//	log.Fatal(`reading root certificate failed`)
	//}
	//
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(rootCert)
	//return &Client{hc: &http.Client{Transport: &http.Transport{
	//	TLSClientConfig: &tls.Config{
	//		RootCAs: caCertPool,
	//	},
	//}}}

	log.Info("initialized an HTTP client")
	return &Client{hc: &http.Client{}}
}

func (c *Client) Send(data []byte, destination any) (res []byte, err error) {
	addr, ok := destination.(string)
	if !ok {
		return nil, urlStringError(destination)
	}

	var status int
	if data == nil {
		res, status, err = c.get(addr)
		if err != nil {
			return nil, sendFailed(addr, http.MethodGet, err)
		}
	} else {
		res, status, err = c.post(addr, data)
		if err != nil {
			return nil, sendFailed(addr, http.MethodPost, err)
		}
	}

	if status != http.StatusOK && status != http.StatusCreated {
		return res, invalidStatusCode(status)
	}

	return res, nil
}

func (c *Client) get(addr string) (response []byte, statusCode int, err error) {
	res, err := c.hc.Get(addr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send HTTP GET request: %w", err)
	}
	defer res.Body.Close()

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	return resData, res.StatusCode, nil
}

func (c *Client) post(url string, data []byte) (resData []byte, statusCode int, err error) {
	res, err := c.hc.Post(url, `application/json`, bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send HTTP POST request: %w", err)
	}
	defer res.Body.Close()

	resData, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	return resData, res.StatusCode, nil
}
