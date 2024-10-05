package transfer

type Request struct {
	TransferFormat   string `json:"transferFormat"`
	AgreementId      string `json:"agreementId"`
	ProviderEndpoint string `json:"providerEndpoint"`
	// integrate token
	DataSink struct {
		Database string `json:"database"`
		Endpoint string `json:"endpoint"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"dataSink"`
}

type StartRequest struct {
	Provider       bool   `json:"provider"`
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

type TerminateRequest struct {
	Provider   bool          `json:"provider"`
	TransferId string        `json:"transferProcessId"`
	Code       string        `json:"code"`
	Reasons    []interface{} `json:"reasons"`
}
