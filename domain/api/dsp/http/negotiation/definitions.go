package negotiation

type State string

// States used during the protocol
const (
	StateRequested  State = "REQUESTED"
	StateOffered    State = "OFFERED"
	StateAccepted   State = "ACCEPTED"
	StateAgreed     State = "AGREED"
	StateVerified   State = "VERIFIED"
	StateFinalized  State = "FINALIZED"
	StateTerminated State = "TERMINATED"
)

type EventType string

// Event types used in the notification message
const (
	EventAccepted  EventType = `dspace:ACCEPTED`
	EventFinalized EventType = `dspace:FINALIZED`
)

// Path parameters
const (
	ParamProviderId  = `providerPid`
	ParamConsumerPid = `consumerPid`
)

// Message types
const (
	MsgTypeNegotiation           = `dspace:ContractNegotiation`
	MsgTypeContractOffer         = `dspace:ContractOfferMessage`
	MsgTypeContractRequest       = `dspace:ContractRequestMessage`
	MsgTypeNegotiationAck        = `dspace:ContractNegotiationAckMessage`
	MsgTypeContractAgreement     = `dspace:ContractAgreementMessage`
	MsgTypeAgreementVerification = `dspace:ContractAgreementVerificationMessage`
	MsgTypeNegotiationEvent      = `dspace:ContractNegotiationEventMessage`
)

// Endpoints
const (
	RequestEndpoint               = `/negotiations/{` + ParamProviderId + `}`
	ContractRequestEndpoint       = `/negotiations/request`
	ContractAgreementEndpoint     = `/negotiations/{` + ParamConsumerPid + `}/agreement`
	AgreementVerificationEndpoint = `/negotiations/{` + ParamProviderId + `}/agreement/verification`
	EventConsumerEndpoint         = `/negotiations/{` + ParamConsumerPid + `}/events`
	EventProviderEndpoint         = `/negotiations/{` + ParamProviderId + `}/events`
)
