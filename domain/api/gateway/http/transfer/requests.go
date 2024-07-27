package transfer

type Request struct {
	TransferFormat   string `json:"transferFormat"`
	AgreementId      string `json:"agreementId"`
	SinkEndpoint     string `json:"sinkEndpoint"`
	ProviderEndpoint string `json:"providerEndpoint"`
}

type StartRequest struct {
	TransferId     string `json:"transferProcessId"`
	SourceEndpoint string `json:"sourceEndpoint"`
}
