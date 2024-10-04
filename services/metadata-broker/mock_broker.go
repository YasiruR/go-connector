package metadata_broker

import (
	"github.com/YasiruR/go-connector/domain/services/metadata-broker"
)

type CeitBroker struct{}

func (c *CeitBroker) SD() (metadata_broker.SelfDescription, error) {
	return metadata_broker.SelfDescription{}, nil
}

func (c *CeitBroker) AddSD(sd metadata_broker.SelfDescription) (id string, err error) {
	return ``, nil
}

func (c *CeitBroker) UpdateSD(id string, sd metadata_broker.SelfDescription) error { return nil }

func (c *CeitBroker) EnableSD(id string) error { return nil }

func (c *CeitBroker) DisableSD(id string) error { return nil }

func (c *CeitBroker) SDByConnector(conId string) (metadata_broker.SelfDescription, error) {
	return metadata_broker.SelfDescription{}, nil
}

func (c *CeitBroker) SDByQuery(query string) (metadata_broker.SelfDescription, error) {
	return metadata_broker.SelfDescription{}, nil
}
