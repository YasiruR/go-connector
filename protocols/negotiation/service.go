package negotiation

type Consumer interface {
	RequestContract()
	AcceptContract()
	VerifyAgreement()
	TerminateContract()
}

type Provider interface {
	OfferContract()
	AgreeContract()
	FinalizeContract()
}
