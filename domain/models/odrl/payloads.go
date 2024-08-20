package odrl

const (
	TypeOffer     = `odrl:Offer`
	TypeAgreement = `odrl:Agreement`
)

// action types
const (
	ActionUse = `odrl:use`
)

type Action string
type Assigner string
type Assignee string
type Target string

// Offer is a subclass of Policy that supports offerings of Rules from assigner Parties. It
// usually does not contain Assignee but DSP defines Assignee in ContractRequest
// (see: https://github.com/International-Data-Spaces-Association/ids-specification/blob/main/negotiation/message/example/contract-offer-message.json).
type Offer struct {
	Id           string   `json:"@id"`
	Type         string   `json:"@type" default:"odrl:Offer"`
	Target       Target   `json:"odrl:target"`
	Assigner     Assigner `json:"odrl:assigner"`
	Assignee     Assignee `json:"odrl:assignee"`
	Permissions  []Rule   `json:"odrl:permission"`
	Prohibitions []Rule   `json:"odrl:prohibition"`
	// duties/obligations should be included
}

// Agreement is a subclass of Policy that supports granting of Rules from assigner to assignee Parties
type Agreement struct {
	Id          string   `json:"@id"`
	Type        string   `json:"@type" default:"odrl:Agreement"`
	Target      Target   `json:"odrl:target"`
	Assigner    Assigner `json:"odrl:assigner"`
	Assignee    Assignee `json:"odrl:assignee"`
	Timestamp   string   `json:"dspace:timestamp"` // due to this attribute may need to transfer agreement structure to dsp api
	Permissions []Rule   `json:"odrl:permission"`
}

type Rule struct {
	Action      Action       `json:"odrl:action"`
	Constraints []Constraint `json:"odrl:constraint"`
}

type Constraint struct {
	LeftOperand  string `json:"odrl:leftOperand"`
	Operator     string `json:"odrl:operator"`
	RightOperand string `json:"odrl:rightOperand"`
	//RightOperand struct {					// support both strings and structs
	//	Value string `json:"@value"`
	//	Type  string `json:"@type"`
	//} `json:"odrl:rightOperand"`
}
