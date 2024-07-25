package transfer

type State string

const (
	StateRequested State = `dspace:REQUESTED`
	StateStarted   State = `dspace:STARTED`
)

type DataTransferType string

const (
	HTTPPull DataTransferType = `HTTP_PULL`
	HTTPPush DataTransferType = `HTTP_PUSH`
)

// Message types
const (
	TypeTransferProcess  = `dspace:TransferProcess`
	TypeTransferRequest  = `dspace:TransferRequestMessage`
	TypeTransferStart    = `dspace:TransferStartMessage`
	TypeDataAddress      = `dspace:DataAddress`
	TypeEndpointProperty = `dspace:EndpointProperty`
)

// Path parameters
const (
	ParamConsumerPid = `consumerPid`
)

// Endpoints
const (
	RequestEndpoint       = `/transfers/request`
	StartTransferEndpoint = `/transfers/{` + ParamConsumerPid + `}/start`
)

const (
	EndpointTypeHTTP = `https://w3id.org/idsa/v4.1/HTTP`
)
