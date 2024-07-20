package negotiation

import "github.com/YasiruR/connector/core/protocols/odrl"

const (
	TypeProviderHandler = `ProviderHandler`
)

type Consumer interface {
	RequestContract(offerId, providerEndpoint, providerPid string, ot odrl.Target, a odrl.Assigner, act odrl.Action) error
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
