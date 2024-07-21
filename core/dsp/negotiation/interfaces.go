package negotiation

const (
	TypeProviderHandler = `ProviderHandler`
)

type Consumer interface {
	RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) error
	AcceptContract()
	VerifyAgreement()
	TerminateContract()
}

type Provider interface {
	OfferContract()
	AgreeContract(offerId string)
	FinalizeContract()
}

type ProviderHandler interface {
	HandleNegotiationsRequest(providerPid string) (Ack, error)
	HandleContractRequest(cr ContractRequest) (Ack, error)
}

type StateMachine interface {
	Requested()
	Offered()
	Accepted()
	Agreed()
	Verified()
	Terminated()
	Finalized()
}
