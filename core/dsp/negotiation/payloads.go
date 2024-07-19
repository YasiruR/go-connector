package negotiation

// Data models required for the negotiation process as specified by
// https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/contract-negotiation/contract.negotiation.protocol

// request payloads used in Negotiation Protocol

type ContractOffer struct {
	Ctx          string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string `json:"@type" default:"dspace:ContractOfferMessage"`
	ProvPId      string `json:"dspace:providerPid"`
	ConsPId      string `json:"dspace:consumerPid"`
	Offer        Offer  `json:"dspace:offer"`
	CallbackAddr string `json:"dspace:callbackAddress"`
}

type ContractRequest struct {
	Ctx          string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string `json:"@type" default:"dspace:ContractRequestMessage"`
	ProvPId      string `json:"dspace:providerPid"`
	ConsPId      string `json:"dspace:consumerPid"`
	Offer        Offer  `json:"dspace:offer"`
	CallbackAddr string `json:"dspace:callbackAddress"`
}

type ContractAgreement struct {
	Ctx          string    `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string    `json:"@type" default:"dspace:ContractAgreementMessage"`
	ProvPId      string    `json:"dspace:providerPid"`
	ConsPId      string    `json:"dspace:consumerPid"`
	Agreement    Agreement `json:"dspace:agreement"`
	CallbackAddr string    `json:"dspace:callbackAddress"`
}

type ContractVerification struct {
	Ctx     string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string `json:"@type" default:"dspace:ContractAgreementVerificationMessage"`
	ProvPId string `json:"dspace:providerPid"`
	ConsPId string `json:"dspace:consumerPid"`
}

type ContractNegotiation struct {
	Ctx       string    `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type      string    `json:"@type" default:"dspace:ContractNegotiationEventMessage"`
	ProvPId   string    `json:"dspace:providerPid"`
	ConsPId   string    `json:"dspace:consumerPid"`
	EventType EventType `json:"dspace:eventType"`
}

type ContractTermination struct {
	Ctx     string   `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string   `json:"@type" default:"dspace:ContractNegotiationTerminationMessage"`
	ProvPId string   `json:"dspace:providerPid"`
	ConsPId string   `json:"dspace:consumerPid"`
	Code    string   `json:"dspace:code"`
	Reason  []Reason `json:"dspace:reason"`
}

// response types used in Negotiation Protocol

type Ack struct {
	Ctx     string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string `json:"@type" default:"dspace:ContractNegotiation"`
	ProvPId string `json:"dspace:providerPid"`
	ConsPId string `json:"dspace:consumerPid"`
	State   State  `json:"dspace:state"`
}

func NewAck() Ack {
	return Ack{
		Ctx:  "https://w3id.org/dspace/2024/1/context.json",
		Type: "dspace:ContractNegotiationAckMessage",
	}
}

type Error struct {
	Ctx     string        `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string        `json:"@type" default:"dspace:ContractNegotiationError"`
	ProvPId string        `json:"dspace:providerPid"`
	ConsPId string        `json:"dspace:consumerPid"`
	Code    string        `json:"dspace:code"`
	Reason  []interface{} `json:"dspace:reason"` // minItems: 1
	Desc    []struct {
		Lang string `json:"@language"`
		Val  string `json:"@value"`
	} `json:"dct:description"`
}

// nested payload structures required by Negotiation Protocol

type Offer struct {
	Id          string       `json:"@id"`
	Type        string       `json:"@type" default:"odrl:Offer"`
	Target      string       `json:"odrl:target"`
	Assigner    string       `json:"odrl:assigner"`
	Permissions []Permission `json:"odrl:permission"`
}

type Agreement struct {
	Id        string `json:"@id"`
	Type      string `json:"@type" default:"odrl:Agreement"`
	Target    string `json:"odrl:target"`
	Assigner  string `json:"odrl:assigner"`
	Assignee  string `json:"odrl:assignee"`
	Timestamp string `json:"dspace:timestamp"`
}

type Permission struct {
	Action      string       `json:"odrl:action"`
	Constraints []Constraint `json:"odrl:constraint"`
}

type Constraint struct {
	LeftOperand  string `json:"odrl:leftOperand"`
	Operand      string `json:"odrl:operand"`
	RightOperand struct {
		Value string `json:"@value"`
		Type  string `json:"@type"`
	} `json:"odrl:rightOperand"`
}

type Reason struct {
	Value    string `json:"@value"`
	Language string `json:"@language"`
}
