package negotiation

type ContractRequest struct {
	ConsumerPId      string            `json:"consumerPid"`
	OfferId          string            `json:"offerId"`
	ProviderEndpoint string            `json:"providerEndpoint"`
	Constraints      map[string]string `json:"constraints"`
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

type TerminateContractRequest struct {
	ConsumerPid string   `json:"consumerPid"`
	ProviderPid string   `json:"providerPid"`
	Code        string   `json:"code"`
	Reasons     []string `json:"reasons"`
}
