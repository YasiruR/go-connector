package transfer

import "github.com/YasiruR/go-connector/domain/api"

const (
	GetProcessEndpoint = `/gateway/transfer/get-process/{` + api.ParamConsumerPid + `}`
	RequestEndpoint    = `/gateway/transfer/request`
	StartEndpoint      = `/gateway/transfer/start`
	SuspendEndpoint    = `/gateway/transfer/suspend`
	CompleteEndpoint   = `/gateway/transfer/complete`
	TerminateEndpoint  = `/gateway/transfer/terminate`
)
