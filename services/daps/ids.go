package daps

import (
	"github.com/tryfix/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type IDS struct {
}

func (i *IDS) Metadata() {}

func (i *IDS) Register(jwt string) {
	data := url.Values{}
	// client credentials grant rfc: https://datatracker.ietf.org/doc/html/rfc6749#section-4.4
	data.Set(`grant_type`, `client_credentials`)
	// use jwt for client authentication rfc: https://datatracker.ietf.org/doc/html/rfc7523#section-2.2
	data.Set(`&client_assertion_type`, `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`)
	data.Set(`&client_assertion`, jwt)
	// set of attributes to be requested from DAPS
	data.Set(`&scope`, `idsc:IDS_CONNECTOR_ATTRIBUTES_ALL`)

	req, err := http.NewRequest(http.MethodPost, `https://localhost:443/token`, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(`new request failed`, err)
	}
	req.Header.Add("Content-Type", `application/x-www-form-urlencoded`)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(`do request failed`, err)
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(`read response body failed`, err)
	}
}
