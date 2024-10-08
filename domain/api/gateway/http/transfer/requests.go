package transfer

import "github.com/YasiruR/go-connector/domain/models"

type Request struct {
	TransferFormat   string `json:"transferFormat"`
	AgreementId      string `json:"agreementId"`
	ProviderEndpoint string `json:"providerEndpoint"`
	// integrate token
	DataSink struct {
		// other sink types
		models.Database
	} `json:"dataSink"`
}

type StartRequest struct {
	Provider   bool   `json:"provider"`
	TransferId string `json:"transferProcessId"`
	DataSource struct {
		models.Database
	} `json:"dataSource"`
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
