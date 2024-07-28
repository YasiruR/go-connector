package transfer

type State string

const (
	StateRequested State = `dspace:REQUESTED`
	StateStarted   State = `dspace:STARTED`
	StateSuspended State = `dspace:SUSPENDED`
	StateCompleted State = `dspace:COMPLETED`
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
	MsgTypeDataAddress      = `dspace:DataAddress`
	MsgTypeEndpointProperty = `dspace:EndpointProperty`
)

// Path parameters
const (
	ParamPid         = `Pid`
	ParamConsumerPid = `consumerPid`
)

// Endpoints
const (
	RequestEndpoint  = `/transfers/request`
	StartEndpoint    = `/transfers/{` + ParamConsumerPid + `}/start`
	SuspendEndpoint  = `/transfers/{` + ParamPid + `}/suspension`
	CompleteEndpoint = `/transfers/{` + ParamPid + `}/completion`
)

const (
	EndpointTypeHTTP = `https://w3id.org/idsa/v4.1/HTTP`
)
