package control_plane

import (
	"github.com/YasiruR/go-connector/domain/control-plane/consumer"
	"github.com/YasiruR/go-connector/domain/control-plane/owner"
	"github.com/YasiruR/go-connector/domain/control-plane/provider"
)

const (
	RoleProvider = `provider`
	RoleConsumer = `consumer`
	RoleOwner    = `owner`
)

type Provider interface {
	provider.CatalogHandler
	provider.NegotiationController
	provider.NegotiationHandler
	provider.TransferController
	provider.TransferHandler
}

type Consumer interface {
	consumer.CatalogController
	consumer.NegotiationController
	consumer.NegotiationHandler
	consumer.TransferController
	consumer.TransferHandler
}

type Owner interface {
	owner.Controller
}
