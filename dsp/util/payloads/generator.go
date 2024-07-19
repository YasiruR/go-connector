package payloads

import "github.com/YasiruR/connector/core/dsp/negotiation"

func NewNegotiation() negotiation.Negotiation {
	return negotiation.Negotiation{
		Ctx:  "https://w3id.org/dspace/2024/1/context.json",
		Type: "dspace:ContractNegotiationAckMessage",
	}
}
