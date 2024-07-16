package negotiation

// Data models required for the negotiation process as specified by
// https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/contract-negotiation/contract.negotiation.protocol

// todo check

type ContractOffer struct {
	Ctx          string   `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string   `json:"@type" default:"dspace:ContractOfferMessage"`
	ProvId       string   `json:"dspace:providerPid"`
	ConsId       string   `json:"dspace:consumerPid"`
	Offer        struct{} `json:"dspace:offer"`
	CallbackAddr string   `json:"dspace:callbackAddress"`
}

type ContractRequest struct {
	Ctx          string   `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string   `json:"@type" default:"dspace:ContractRequestMessage"`
	ProvId       string   `json:"dspace:providerPid"`
	ConsId       string   `json:"dspace:consumerPid"`
	Offer        struct{} `json:"dspace:offer"`
	CallbackAddr string   `json:"dspace:callbackAddress"`
}

type ContractAgreement struct {
	Ctx          string   `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string   `json:"@type" default:"dspace:ContractAgreementMessage"`
	ProvId       string   `json:"dspace:providerPid"`
	ConsId       string   `json:"dspace:consumerPid"`
	Agreement    struct{} `json:"dspace:agreement"`
	CallbackAddr string   `json:"dspace:callbackAddress"`
}

type ContractVerification struct {
	Ctx    string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type   string `json:"@type" default:"dspace:ContractAgreementVerificationMessage"`
	ProvId string `json:"dspace:providerPid"`
	ConsId string `json:"dspace:consumerPid"`
}

type ContractNegotiation struct {
	Ctx       string    `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type      string    `json:"@type" default:"dspace:ContractNegotiationEventMessage"`
	ProvId    string    `json:"dspace:providerPid"`
	ConsId    string    `json:"dspace:consumerPid"`
	EventType EventType `json:"dspace:eventType"`
}

type ContractTermination struct {
	Ctx    string        `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type   string        `json:"@type" default:"dspace:ContractNegotiationEventMessage"`
	ProvId string        `json:"dspace:providerPid"`
	ConsId string        `json:"dspace:consumerPid"`
	Code   string        `json:"dspace:code"`
	Reason []interface{} `json:"dspace:reason"` // minItems: 1
}

type ContractNegotiationAck struct {
	Ctx    string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type   string `json:"@type" default:"dspace:ContractNegotiationEventMessage"`
	ProvId string `json:"dspace:providerPid"`
	ConsId string `json:"dspace:consumerPid"`
	State  State  `json:"dspace:state"`
}

type ContractNegotiationErr struct {
	Ctx    string        `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type   string        `json:"@type" default:"dspace:ContractNegotiationEventMessage"`
	ProvId string        `json:"dspace:providerPid"`
	ConsId string        `json:"dspace:consumerPid"`
	Code   string        `json:"dspace:code"`
	Reason []interface{} `json:"dspace:reason"` // minItems: 1
	Desc   []struct {
		Lang string `json:"@language"`
		Val  string `json:"@value"`
	} `json:"dct:description"`
}
