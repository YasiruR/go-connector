package transfer

type Request struct {
	Ctx          string           `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string           `json:"@type" default:"dspace:TransferRequestMessage"`
	ConsPId      string           `json:"dspace:consumerPid"`
	AgreementId  string           `json:"dspace:agreementId"`
	Format       DataTransferType `json:"dct:format"`
	Address      Address          `json:"dspace:address"` // required only if format is a push transfer
	CallbackAddr string           `json:"dspace:callbackAddress"`
}

type StartRequest struct {
	Ctx     string  `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string  `json:"@type" default:"dspace:TransferStartMessage"`
	ConsPId string  `json:"dspace:consumerPid"`
	ProvPId string  `json:"dspace:providerPid"`
	Address Address `json:"dspace:address"`
}

type SuspendRequest struct {
	Ctx     string        `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string        `json:"@type" default:"dspace:TransferSuspensionMessage"`
	ConsPId string        `json:"dspace:consumerPid"`
	ProvPId string        `json:"dspace:providerPid"`
	Code    string        `json:"dspace:code"`
	Reason  []interface{} `json:"dspace:reason"`
}

type CompleteRequest struct {
	Ctx     string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string `json:"@type" default:"dspace:TransferCompletionMessage"`
	ConsPId string `json:"dspace:consumerPid"`
	ProvPId string `json:"dspace:providerPid"`
}

type TerminateRequest struct {
	Ctx     string        `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type    string        `json:"@type" default:"dspace:TransferTerminationMessage"`
	ConsPId string        `json:"dspace:consumerPid"`
	ProvPId string        `json:"dspace:providerPid"`
	Code    string        `json:"dspace:code"`
	Reason  []interface{} `json:"dspace:reason"`
}

type Address struct {
	Type               string             `json:"@type" default:"dspace:DataAddress"`
	EndpointType       string             `json:"dspace:endpointType"`
	Endpoint           string             `json:"dspace:endpoint"`
	EndpointProperties []EndpointProperty `json:"dspace:endpointProperties"` // include authorization details for the consumer endpoint (e.g. token)
}

type EndpointProperty struct {
	Type  string `json:"@type" default:"dspace:EndpointProperty"`
	Name  string `json:"dspace:name"`
	Value string `json:"dspace:value"`
}
