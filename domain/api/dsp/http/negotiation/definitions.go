package negotiation

import "github.com/YasiruR/connector/domain/api"

type State string

// States used during the negotiation protocol
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

// Message types
const (
	MsgTypeNegotiation           = `dspace:ContractNegotiation`
	MsgTypeContractOffer         = `dspace:ContractOfferMessage`
	MsgTypeContractRequest       = `dspace:ContractRequestMessage`
	MsgTypeNegotiationAck        = `dspace:ContractNegotiationAckMessage`
	MsgTypeContractAgreement     = `dspace:ContractAgreementMessage`
	MsgTypeAgreementVerification = `dspace:ContractAgreementVerificationMessage`
	MsgTypeNegotiationEvent      = `dspace:ContractNegotiationEventMessage`
	MsgTypeTermination           = `dspace:ContractNegotiationTerminationMessage`
)

// Provider endpoints
const (
	RequestEndpoint                = `/negotiations/{` + api.ParamProviderPid + `}`
	ContractRequestEndpoint        = `/negotiations/request`
	ContractRequestToOfferEndpoint = `/negotiations/{` + api.ParamProviderPid + `}/request`
	AgreementVerificationEndpoint  = `/negotiations/{` + api.ParamProviderPid + `}/agreement/verification`
)

// Consumer endpoints
const (
	ContractOfferEndpoint          = `/negotiations/offers`
	ContractOfferToRequestEndpoint = `/negotiations/{` + api.ParamConsumerPid + `}/offers`
	ContractAgreementEndpoint      = `/negotiations/{` + api.ParamConsumerPid + `}/agreement`
)

// Common endpoints
const (
	EventsEndpoint    = `/negotiations/{` + api.ParamPid + `}/events`
	TerminateEndpoint = `/negotiations/{` + api.ParamPid + `}/termination`
)
