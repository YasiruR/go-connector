package transfer

type Request struct {
	TransferType     string `json:"transferType"`
	AgreementId      string `json:"agreementId"`
	SinkEndpoint     string `json:"sinkEndpoint"`
	ProviderEndpoint string `json:"providerEndpoint"`
}
