package negotiation

type ControllerProvider interface {
	OfferContract()
	AgreeContract(offerId, negotiationId string) (agreementId string, err error)
	FinalizeContract(providerPid string) error
}

type HandlerProvider interface {
	HandleNegotiationsRequest(providerPid string) (Ack, error)
	HandleContractRequest(cr ContractRequest) (Ack, error)
	HandleAgreementVerification(providerPid string) (Ack, error)
}

type ControllerConsumer interface {
	// change endpoint to generic
	RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) (cnId string, err error)
	AcceptContract()
	VerifyAgreement(consumerPid string) error
	TerminateContract()
}

type HandlerConsumer interface {
	HandleContractAgreement(ca ContractAgreement) (Ack, error)
	HandleFinalizedEvent(consumerPid string) (Ack, error)
}
