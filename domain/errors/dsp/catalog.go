package dsp

import (
	"fmt"
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/core"
)

type CatalogError struct {
	err  error
	Body catalog.Error
}

func (c CatalogError) Error() string {
	return c.err.Error()
}

func NewCatalogError(body catalog.Error, err error) error {
	return CatalogError{
		err:  err,
		Body: body,
	}
}

func CatalogInvalidKey(store, key string, err error) error {
	return CatalogError{
		err: fmt.Errorf("%s store error - %s for %s", store, err, key),
		Body: catalog.Error{
			Context:    core.Context,
			Type:       catalog.MsgTypeError,
			DspaceCode: "2001",
			DspaceReason: []struct {
				Value    string `json:"@value"`
				Language string `json:"@language"`
			}{{Value: fmt.Sprintf("incorrect value provided for %s", key), Language: "en"}},
		},
	}
}

func CatalogInvalidReqBody(msgType string, err error) error {
	return CatalogError{
		err: fmt.Errorf("unmarshal failed for request %s - %w", msgType, err),
		Body: catalog.Error{
			Context:    core.Context,
			Type:       catalog.MsgTypeError,
			DspaceCode: "2002",
			DspaceReason: []struct {
				Value    string `json:"@value"`
				Language string `json:"@language"`
			}{{Value: fmt.Sprintf("invalid request body for %s", msgType), Language: "en"}},
		},
	}
}

func CatalogReqParseError(msgType string, err error) error {
	return CatalogError{
		err: fmt.Errorf("reading request body failed for %s - %s", msgType, err),
		Body: catalog.Error{
			Context:    core.Context,
			Type:       catalog.MsgTypeError,
			DspaceCode: "2003",
			DspaceReason: []struct {
				Value    string `json:"@value"`
				Language string `json:"@language"`
			}{{Value: "request parser error", Language: "en"}},
		},
	}
}

func CatalogWriteAckError(msgType string, err error) error {
	return CatalogError{
		err: fmt.Errorf("%s for %s", err, msgType),
		Body: catalog.Error{
			Context:    core.Context,
			Type:       catalog.MsgTypeError,
			DspaceCode: "2004",
			DspaceReason: []struct {
				Value    string `json:"@value"`
				Language string `json:"@language"`
			}{{Value: "internal error", Language: "en"}},
		},
	}
}
