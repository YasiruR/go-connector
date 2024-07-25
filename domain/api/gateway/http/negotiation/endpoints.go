package negotiation

// Path parameters
const (
	ParamAgreementId = `agreementId`
	ParamConsumerPid = `consumerPid`
	ParamProviderPid = `providerPid`
)

// endpoints exposed by gateway API
const (
	RequestContractEndpoint  = `/gateway/contract`
	AgreeContractEndpoint    = `/gateway/agree-contract`
	GetAgreementEndpoint     = `/gateway/agreement/{` + ParamAgreementId + `}`
	VerifyAgreementEndpoint  = `/gateway/verify-agreement/{` + ParamConsumerPid + `}`
	FinalizeContractEndpoint = `/gateway/finalize-contract/{` + ParamProviderPid + `}`
)
