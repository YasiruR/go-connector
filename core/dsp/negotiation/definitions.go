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
)

type EventType string

const (
	Accepted  EventType = `ACCEPTED`
	Finalized EventType = `FINALIZED`
)

const (
	RequestEndpoint               = `/negotiations/{` + ParamProviderId + `}`
	ContractRequestEndpoint       = `/negotiations/request`
	ContractAgreementEndpoint     = `/negotiations/{` + ParamConsumerPid + `}/agreement`
	AgreementVerificationEndpoint = `/negotiations/{` + ParamProviderId + `}/agreement/verification`
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
