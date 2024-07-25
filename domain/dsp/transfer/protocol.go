package transfer

type Controller interface {
	RequestTransfer(typ DataTransferType, agreementId, sinkEndpoint, providerEndpoint string) (tpId string, err error)
}
