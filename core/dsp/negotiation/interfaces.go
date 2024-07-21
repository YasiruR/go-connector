package negotiation

const (
	TypeProviderHandler = `ProviderHandler`
	TypeConsumerHandler = `ConsumerHandler`
)

type Consumer interface {
	RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) (negotiationId string, err error)
	AcceptContract()
	VerifyAgreement()
	TerminateContract()
}

type ConsumerHandler interface {
	HandleContractAgreement(consumerPid string, ca ContractAgreement) (Ack, error)
}

type Provider interface {
	OfferContract()
	AgreeContract(offerId, negotiationId string) (agreementId string, err error)
	FinalizeContract()
}

type ProviderHandler interface {
	HandleNegotiationsRequest(providerPid string) (Ack, error)
	HandleContractRequest(cr ContractRequest) (Ack, error)
}
