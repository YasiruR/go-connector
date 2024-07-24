package negotiation

type Provider interface {
	ProviderController
	ProviderHandler
}

type ProviderController interface {
	OfferContract()
	AgreeContract(offerId, negotiationId string) (agreementId string, err error)
	FinalizeContract(providerPid string) error
}

type ProviderHandler interface {
	HandleNegotiationsRequest(providerPid string) (Ack, error)
	HandleContractRequest(cr ContractRequest) (Ack, error)
	HandleAgreementVerification(providerPid string) (Ack, error)
}
