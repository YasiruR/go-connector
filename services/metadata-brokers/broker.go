package metadata_brokers

import "github.com/YasiruR/connector/domain/services"

type CeitBroker struct{}

func (c *CeitBroker) SD() (services.SelfDescription, error) { return services.SelfDescription{}, nil }

func (c *CeitBroker) AddSD(sd services.SelfDescription) (id string, err error) { return ``, nil }

func (c *CeitBroker) UpdateSD(id string, sd services.SelfDescription) error { return nil }

func (c *CeitBroker) EnableSD(id string) error { return nil }

func (c *CeitBroker) DisableSD(id string) error { return nil }

func (c *CeitBroker) SDByConnector(conId string) (services.SelfDescription, error) {
	return services.SelfDescription{}, nil
}

func (c *CeitBroker) SDByQuery(query string) (services.SelfDescription, error) {
	return services.SelfDescription{}, nil
}
