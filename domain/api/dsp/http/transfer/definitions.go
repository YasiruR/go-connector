package transfer

import "github.com/YasiruR/go-connector/domain/api"

type State string

const (
	StateRequested  State = `dspace:REQUESTED`
	StateStarted    State = `dspace:STARTED`
	StateSuspended  State = `dspace:SUSPENDED`
	StateCompleted  State = `dspace:COMPLETED`
	StateTerminated State = `dspace:TERMINATED`
)

type DataTransferType string

const (
	HTTPPull DataTransferType = `HTTP_PULL`
	HTTPPush DataTransferType = `HTTP_PUSH`
)

// Message types
const (
	MsgTypeProcess          = `dspace:TransferProcess`
	MsgTypeRequest          = `dspace:TransferRequestMessage`
	MsgTypeStart            = `dspace:TransferStartMessage`
	MsgTypeSuspend          = `dspace:TransferSuspensionMessage`
	MsgTypeComplete         = `dspace:TransferCompletionMessage`
	MsgTypeTerminate        = `dspace:TransferTerminationMessage`
	MsgTypeError            = `dspace:TransferError`
	MsgTypeDataAddress      = `dspace:DataAddress`
	MsgTypeEndpointProperty = `dspace:EndpointProperty`
)

// Provider endpoints
const (
	GetProcessEndpoint = `/transfers/{` + api.ParamProviderPid + `}`
	RequestEndpoint    = `/transfers/request`
)

// Common endpoints
const (
	StartEndpoint     = `/transfers/{` + api.ParamPid + `}/start`
	SuspendEndpoint   = `/transfers/{` + api.ParamPid + `}/suspension`
	CompleteEndpoint  = `/transfers/{` + api.ParamPid + `}/completion`
	TerminateEndpoint = `/transfers/{` + api.ParamPid + `}/termination`
	EndpointTypeHTTP  = `https://w3id.org/idsa/v4.1/HTTP`
)
