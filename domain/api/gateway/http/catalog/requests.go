package catalog

type Request struct {
	ProviderEndpoint string `json:"providerEndpoint"`
}

type DatasetRequest struct {
	DatasetId        string `json:"datasetId"`
	ProviderEndpoint string `json:"providerEndpoint"`
}

type CreatePolicyRequest struct {
	Permissions  []Rule `json:"permissions"`
	Prohibitions []Rule `json:"prohibitions"`
	Obligations  []Rule `json:"obligations"`
}

type CreateDatasetRequest struct {
	Title        string   `json:"title"`
	Descriptions []string `json:"descriptions"`
	Endpoints    []string `json:"endpoints"`
	OfferIds     []string `json:"offerIds"`
	Keywords     []string `json:"keywords"`
	Format       string   `json:"format"`
}

type Rule struct {
	Action      string       `json:"action"`
	Constraints []Constraint `json:"constraints"`
}

type Constraint struct {
	LeftOperand  string `json:"leftOperand"`
	Operator     string `json:"operator"`
	RightOperand string `json:"rightOperand"`
}
