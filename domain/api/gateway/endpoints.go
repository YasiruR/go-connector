package gateway

const (
	ParamAgreementId = `agreementId`
	ParamConsumerPid = `consumerPid`
	ParamProviderPid = `providerPid`
)

// endpoints exposed by gateway API
const (
	CreatePolicyEndpoint     = `/gateway/create-policy`
	CreateDatasetEndpoint    = `/gateway/create-dataset`
	RequestCatalogEndpoint   = `/gateway/catalog`
	RequestDatasetEndpoint   = `/gateway/dataset`
	RequestContractEndpoint  = `/gateway/contract`
	AgreeContractEndpoint    = `/gateway/agree-contract`
	GetAgreementEndpoint     = `/gateway/agreement/{` + ParamAgreementId + `}`
	VerifyAgreementEndpoint  = `/gateway/verify-agreement/{` + ParamConsumerPid + `}`
	FinalizeContractEndpoint = `/gateway/finalize-contract/{` + ParamProviderPid + `}`
)
