package api

import (
	"strings"
)

type Server interface {
	Start()
}

// Query parameters

const (
	ParamPid         = `pid`
	ParamConsumerPid = `consumerPid`
	ParamProviderPid = `providerPid`
)

func SetParamConsumerPid(endpoint, id string) string {
	return strings.Replace(endpoint, `{`+ParamConsumerPid+`}`, id, 1)
}

func SetParamProviderPid(endpoint, id string) string {
	return strings.Replace(endpoint, `{`+ParamProviderPid+`}`, id, 1)
}

func SetParamPid(endpoint, id string) string {
	return strings.Replace(endpoint, `{`+ParamPid+`}`, id, 1)
}
