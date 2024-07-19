package gateway

type PolicyRequest struct {
	Permissions  []Rule `json:"permissions"`
	Prohibitions []Rule `json:"prohibitions"`
	Obligations  []Rule `json:"obligations"`
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

type DatasetRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	PolicyId    string   `json:"policyId"`
}

type ContractRequest struct {
	OfferId          string `json:"offerId"`
	ProviderEndpoint string `json:"providerEndpoint"`
	ProviderPId      string `json:"providerPId"`
	OdrlTarget       string `json:"odrlTarget"`
	Assigner         string `json:"assigner"`
	Action           string `json:"action"`
}

type ContractAgreementRequest struct {
}
