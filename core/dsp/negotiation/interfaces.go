package negotiation

type Consumer interface {
	RequestContract(offerId, providerEndpoint, providerPid, odrlTarget, assigner, action string) error
	AcceptContract()
	VerifyAgreement()
	TerminateContract()
}

type Provider interface {
	OfferContract()
	AgreeContract()
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
