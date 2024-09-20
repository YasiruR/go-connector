package transfer

type Ack Process

type Process struct {
	Ctx     string           `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    DataTransferType `json:"@type" default:"dspace:TransferProcess"`
	ProvPId string           `json:"dspace:providerPid"`
	ConsPId string           `json:"dspace:consumerPid"`
	State   State            `json:"dspace:state"`
}

type Error struct {
	Ctx     string        `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string        `json:"@type" default:"dspace:TransferError"`
	ConsPId string        `json:"dspace:consumerPid"`
	ProvPId string        `json:"dspace:providerPid"`
	Code    string        `json:"dspace:code"`
	Reason  []interface{} `json:"dspace:reason"`
}
