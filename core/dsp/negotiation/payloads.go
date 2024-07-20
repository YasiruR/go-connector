package negotiation

import "github.com/YasiruR/connector/core/protocols/odrl"

// Data models required for the negotiation process as specified by
// https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/contract-negotiation/contract.negotiation.protocol

type ProviderPid string
type ConsumerPid string

// request payloads used in Negotiation Protocol

type ContractOffer struct {
	Ctx          string     `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string     `json:"@type" default:"dspace:ContractOfferMessage"`
	ProvPId      string     `json:"dspace:providerPid"`
	ConsPId      string     `json:"dspace:consumerPid"`
	Offer        odrl.Offer `json:"dspace:offer"`
	CallbackAddr string     `json:"dspace:callbackAddress"`
}

type ContractRequest struct {
	Ctx          string     `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string     `json:"@type" default:"dspace:ContractRequestMessage"`
	ProvPId      string     `json:"dspace:providerPid"`
	ConsPId      string     `json:"dspace:consumerPid"`
	Offer        odrl.Offer `json:"dspace:offer"`
	CallbackAddr string     `json:"dspace:callbackAddress"`
}

type ContractAgreement struct {
	Ctx          string         `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string         `json:"@type" default:"dspace:ContractAgreementMessage"`
	ProvPId      string         `json:"dspace:providerPid"`
	ConsPId      string         `json:"dspace:consumerPid"`
	Agreement    odrl.Agreement `json:"dspace:agreement"`
	CallbackAddr string         `json:"dspace:callbackAddress"`
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

type Ack Negotiation

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

type Negotiation struct {
	Ctx     string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string `json:"@type" default:"dspace:ContractNegotiation"`
	ProvPId string `json:"dspace:providerPid"`
	ConsPId string `json:"dspace:consumerPid"`
	State   State  `json:"dspace:state"`
}

type Reason struct {
	Value    string `json:"@value"`
	Language string `json:"@language"`
}
