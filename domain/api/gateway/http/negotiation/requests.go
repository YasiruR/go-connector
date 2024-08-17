package negotiation

type ContractRequest struct {
	OfferId          string `json:"offerId"`
	ConsumerPId      string `json:"consumerPid"`
	ProviderEndpoint string `json:"providerEndpoint"`
	OdrlTarget       string `json:"odrlTarget"`
	Assigner         string `json:"assigner"`
	Assignee         string `json:"assignee"`
	Action           string `json:"action"`
}

type OfferRequest struct {
	ProviderPid  string `json:"providerPid"`
	OfferId      string `json:"offerId"`
	ConsumerAddr string `json:"consumerAddr"`
}

type AgreeContractRequest struct {
	OfferId       string `json:"offerId"`
	NegotiationId string `json:"contractNegotiationId"`
}

type VerifyAgreementRequest struct {
	ConsumerPid string `json:"consumerPid"`
}
