package errors

type CatalogError struct {
	Code string
	ErrorMessage
}

func (c CatalogError) Error() string {
	return c.err.Error()
}

func Catalog(e ErrorMessage) error {
	return CatalogError{
		Code:         "ce_" + e.code,
		ErrorMessage: e,
	}
}

type NegotiationError struct {
	Code    string
	ProvPid string
	ConsPid string
	ErrorMessage
}

func (n NegotiationError) Error() string {
	return n.err.Error()
}

func Negotiation(provPid, consPid string, e ErrorMessage) error {
	return NegotiationError{
		Code:         "ne_" + e.code,
		ProvPid:      provPid,
		ConsPid:      consPid,
		ErrorMessage: e,
	}
}

type TransferError struct {
	Code    string
	ProvPid string
	ConsPid string
	ErrorMessage
}

func (t TransferError) Error() string {
	return t.err.Error()
}

func Transfer(provPid, consPid string, e ErrorMessage) error {
	return TransferError{
		Code:         "te_" + e.code,
		ProvPid:      provPid,
		ConsPid:      consPid,
		ErrorMessage: e,
	}
}

type ClientError struct {
	Code string `json:"code"`
	ErrorMessage
}

func (c ClientError) Error() string {
	return c.err.Error()
}

func Client(e ErrorMessage) error {
	return ClientError{
		Code:         `cl_` + e.code,
		ErrorMessage: e,
	}
}
