package api

import (
	"strings"
)

type Server interface {
	Start()
}

// Query parameters

// todo: revamp negotiation controller and handler (e.g. send func)

const (
	ParamPid         = `pid`
	ParamConsumerPid = `consumerPid`
	ParamProviderPid = `providerPid`
)

func SetConsumerPidParam(endpoint, id string) string {
	return strings.Replace(endpoint, `{`+ParamConsumerPid+`}`, id, 1)
}

func SetProviderPidParam(endpoint, id string) string {
	return strings.Replace(endpoint, `{`+ParamProviderPid+`}`, id, 1)
}
