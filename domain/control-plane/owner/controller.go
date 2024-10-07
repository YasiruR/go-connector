package owner

import "github.com/YasiruR/go-connector/domain/models/odrl"

type Controller interface {
	CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (id string, err error)
	CreateDataset(title, format string, descriptions, keywords, endpoints, offerIds []string) (id string, err error)
}
