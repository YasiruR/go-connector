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

type SuspendRequest struct {
	Provider   bool          `json:"provider"`
	TransferId string        `json:"transferProcessId"`
	Code       string        `json:"code"`
	Reasons    []interface{} `json:"reasons"`
}

type CompleteRequest struct {
	Provider   bool   `json:"provider"`
	TransferId string `json:"transferProcessId"`
}
