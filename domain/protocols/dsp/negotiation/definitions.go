package negotiation

const (
	ParamProviderId  = `providerPid`
	ParamConsumerPid = `consumerPid`
)

const (
	TypeNegotiation           = `dspace:ContractNegotiation`
	TypeContractOffer         = `dspace:ContractOfferMessage`
	TypeContractRequest       = `dspace:ContractRequestMessage`
	TypeNegotiationAck        = `dspace:ContractNegotiationAckMessage`
	TypeContractAgreement     = `dspace:ContractAgreementMessage`
	TypeAgreementVerification = `dspace:ContractAgreementVerificationMessage`
	TypeNegotiationEvent      = `dspace:ContractNegotiationEventMessage`
)

type EventType string

const (
	EventAccepted  EventType = `dspace:ACCEPTED`
	EventFinalized EventType = `dspace:FINALIZED`
)

const (
	RequestEndpoint               = `/negotiations/{` + ParamProviderId + `}`
	ContractRequestEndpoint       = `/negotiations/request`
	ContractAgreementEndpoint     = `/negotiations/{` + ParamConsumerPid + `}/agreement`
	AgreementVerificationEndpoint = `/negotiations/{` + ParamProviderId + `}/agreement/verification`
	EventConsumerEndpoint         = `/negotiations/{` + ParamConsumerPid + `}/events`
	EventProviderEndpoint         = `/negotiations/{` + ParamProviderId + `}/events`
)

type State string

const (
	StateRequested  State = "REQUESTED"
	StateOffered    State = "OFFERED"
	StateAccepted   State = "ACCEPTED"
	StateAgreed     State = "AGREED"
	StateVerified   State = "VERIFIED"
	StateFinalized  State = "FINALIZED"
	StateTerminated State = "TERMINATED"
)
