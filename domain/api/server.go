package api

type HandlerType string

const (
	HandlerTypeCatalog     HandlerType = `catalog`
	HandlerTypeNegotiation HandlerType = `negotiation`
	HandlerTypeTransfer    HandlerType = `transfer`
)

type Server interface {
	Start()
}
