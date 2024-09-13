package services

type SelfDescription struct{}

type MetadataBroker interface {
	SD() (SelfDescription, error)
	AddSD(sd SelfDescription) (id string, err error)
	UpdateSD(id string, sd SelfDescription) error
	EnableSD(id string) error
	DisableSD(id string) error
	SDByConnector(conId string) (SelfDescription, error)
	SDByQuery(query string) (SelfDescription, error)
}
