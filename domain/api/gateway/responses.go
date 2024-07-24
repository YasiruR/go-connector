package gateway

type PolicyResponse struct {
	Id string `json:"policyId"`
}

type DatasetResponse struct {
	Id string `json:"datasetId"`
}

type ContractRequestResponse struct {
	Id string `json:"contractNegotiationId"`
}

type ContractAgreementResponse struct {
	Id string `json:"contractAgreementId"`
}
