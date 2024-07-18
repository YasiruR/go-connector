package negotiation

type EventType string

const (
	Accepted  EventType = `ACCEPTED`
	Finalized EventType = `FINALIZED`
)

const (
	RequestContractEndpoint = `/negotiations/request`
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
