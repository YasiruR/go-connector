package negotiation

type ConsumerController interface {
	// change endpoint to generic
	RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) (negotiationId string, err error)
	AcceptContract()
	VerifyAgreement(consumerPid string) error
	TerminateContract()
}

type ConsumerHandler interface {
	HandleContractAgreement(consumerPid string, ca ContractAgreement) (Ack, error)
}

type ProviderController interface {
	OfferContract()
	AgreeContract(offerId, negotiationId string) (agreementId string, err error)
	FinalizeContract()
}

type ProviderHandler interface {
	HandleNegotiationsRequest(providerPid string) (Ack, error)
	HandleContractRequest(cr ContractRequest) (Ack, error)
	HandleAgreementVerification(providerPid string) (Ack, error)
}
