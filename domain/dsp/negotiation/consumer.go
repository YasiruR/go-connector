package negotiation

type Consumer interface {
	ConsumerController
	ConsumerHandler
}

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
