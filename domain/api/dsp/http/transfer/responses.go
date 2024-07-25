package transfer

type Ack Process

type Process struct {
	Ctx     string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string `json:"@type" default:"dspace:TransferProcess"`
	ProvPId string `json:"dspace:providerPid"`
	ConsPId string `json:"dspace:consumerPid"`
	State   State  `json:"dspace:state"`
}
