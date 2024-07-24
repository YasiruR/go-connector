package negotiation

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
