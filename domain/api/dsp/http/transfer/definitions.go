package transfer

type State string

const (
	StateRequested State = `dspace:REQUESTED`
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
	TypeDataAddress      = `dspace:DataAddress`
	TypeEndpointProperty = `dspace:EndpointProperty`
)

const (
	RequestEndpoint = `/transfers/request`
)

const (
	EndpointTypeHTTP = `https://w3id.org/idsa/v4.1/HTTP`
)
