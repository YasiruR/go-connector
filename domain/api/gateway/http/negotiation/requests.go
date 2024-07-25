package negotiation

type ContractRequest struct {
	OfferId          string `json:"offerId"`
	ProviderEndpoint string `json:"providerEndpoint"`
	OdrlTarget       string `json:"odrlTarget"`
	Assigner         string `json:"assigner"`
	Assignee         string `json:"assignee"`
	Action           string `json:"action"`
}

type AgreeContractRequest struct {
	OfferId       string `json:"offerId"`
	NegotiationId string `json:"negotiationId"`
}

type VerifyAgreementRequest struct {
	ConsumerPid string `json:"consumerPid"`
}
