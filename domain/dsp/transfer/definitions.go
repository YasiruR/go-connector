package transfer

type DataTransferType string

const (
	HTTPPull DataTransferType = `HTTP_PULL`
	HTTPPush DataTransferType = `HTTP_PUSH`
)

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
	StateRequested = `dspace:REQUESTED`
)

const (
	EndpointTypeHTTP = `https://w3id.org/idsa/v4.1/HTTP`
)
