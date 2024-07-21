package negotiation

const (
	Context               = `https://w3id.org/dspace/2024/1/context.json`
	TypeNegotiation       = `dspace:ContractNegotiation`
	TypeOffer             = `dspace:ContractOfferMessage`
	TypeNegotiationAck    = `dspace:ContractNegotiationAckMessage`
	TypeContractAgreement = `dspace:ContractAgreementMessage`
)

type EventType string

const (
	Accepted  EventType = `ACCEPTED`
	Finalized EventType = `FINALIZED`
)

const (
	NegotiationsEndpoint      = `/negotiations/{providerPid}`
	RequestContractEndpoint   = `/negotiations/request`
	ContractAgreementEndpoint = `/negotiations/{consumerPid}/agreement`
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
