package gateway

type ContractRequest struct {
	OfferId          string `json:"offerId"`
	ProviderEndpoint string `json:"providerEndpoint"`
	ProviderPId      string `json:"providerPId"`
	OdrlTarget       string `json:"odrlTarget"`
	Assigner         string `json:"assigner"`
	Action           string `json:"action"`
}

type ContractAgreement struct {
}
