package transfer

type State string

const (
	StateRequested State = `dspace:REQUESTED`
	StateStarted   State = `dspace:STARTED`
	StateSuspended State = `dspace:SUSPENDED`
)

type DataTransferType string

const (
	HTTPPull DataTransferType = `HTTP_PULL`
	HTTPPush DataTransferType = `HTTP_PUSH`
)

// Message types
const (
	MsgTypeProcess         = `dspace:TransferProcess`
	MsgTypeRequest         = `dspace:TransferRequestMessage`
	MsgTypStart            = `dspace:TransferStartMessage`
	MsgTypSuspend          = `dspace:TransferSuspensionMessage`
	MsgTypDataAddress      = `dspace:DataAddress`
	MsgTypEndpointProperty = `dspace:EndpointProperty`
)

// Path parameters
const (
	ParamPid         = `Pid`
	ParamConsumerPid = `consumerPid`
)

// Endpoints
const (
	RequestEndpoint = `/transfers/request`
	StartEndpoint   = `/transfers/{` + ParamConsumerPid + `}/start`
	SuspendEndpoint = `/transfers/{` + ParamPid + `}/suspend`
)

const (
	EndpointTypeHTTP = `https://w3id.org/idsa/v4.1/HTTP`
)
