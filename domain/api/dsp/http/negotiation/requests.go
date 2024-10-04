package negotiation

import (
	"github.com/YasiruR/go-connector/domain/models/odrl"
)

// Data models required for the negotiation process as specified by
// https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/contract-negotiation/contract.negotiation.protocol

// Payloads defined here are the message types that trigger or notify events involved in
// the Negotiation Protocol

type ProviderPid string
type ConsumerPid string

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

type ContractNegotiationEvent struct {
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

type Reason struct {
	Value    string `json:"@value"`
	Language string `json:"@language"`
}
