package gateway

const (
	ParamId          = `id`
	ParamConsumerPid = `consumerPid`
)

// endpoints exposed by gateway API
const (
	CreatePolicyEndpoint    = `/gateway/create-policy`
	CreateDatasetEndpoint   = `/gateway/create-dataset`
	RequestCatalogEndpoint  = `/gateway/catalog`
	RequestDatasetEndpoint  = `/gateway/dataset`
	RequestContractEndpoint = `/gateway/contract`
	AgreeContractEndpoint   = `/gateway/agree-contract`
	GetAgreementEndpoint    = `/gateway/agreement/{` + ParamId + `}`
	VerifyAgreementEndpoint = `/gateway/verify-agreement/{` + ParamConsumerPid + `}`
)
