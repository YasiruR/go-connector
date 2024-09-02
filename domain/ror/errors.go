package ror

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
)

type ClientError struct {
	err  error
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c ClientError) Error() string {
	return c.err.Error()
}

type CatalogError struct {
	err  error
	Body catalog.Error
}

func (c CatalogError) Error() string {
	return c.err.Error()
}

type NegotiationError struct {
	err  error
	Body negotiation.Error
}

func (n NegotiationError) Error() string {
	return n.err.Error()
}

type TransferError struct {
	err  error
	Body transfer.Error
}

func (t TransferError) Error() string {
	return t.err.Error()
}

// internal errors - 1xxx
// protocol errors - 2xxx
// transport failures - 3xxx
