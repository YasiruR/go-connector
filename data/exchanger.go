package data

import (
	"fmt"
	"github.com/YasiruR/go-connector/data/postgresql"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/data"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

type Exchanger struct {
	psql *postgresql.Exchanger
	cat  stores.ProviderCatalog
	log  pkg.Log
}

func NewExchanger(cfg boot.Config, s domain.Stores, log pkg.Log) *Exchanger {
	return &Exchanger{
		cat:  s.ProviderCatalog,
		psql: postgresql.NewExchanger(cfg),
		log:  log,
	}
}

func (e *Exchanger) NewToken(participantId, datasetId string) string {
	return ``
}

func (e *Exchanger) Push(datasetId, host, db, usr, pw string) error {
	return e.push(datasetId, host, db, usr, pw)
}

func (e *Exchanger) Pull(datasetId, endpoint, token string) {

}

func (e *Exchanger) push(datasetId, host, db, usr, pw string) error {
	ds, err := e.cat.Dataset(datasetId)
	if err != nil {
		return errors.Catalog(errors.InvalidKey(stores.TypeProviderCatalog, `dataset id`, err))
	}

	for _, dist := range ds.DcatDistribution {
		switch dist.DctFormat {
		case data.PostgreSQLPush:
			if err = e.psql.Migrate(host, db, usr, pw); err != nil {
				return fmt.Errorf("postgresql failed - %s", err)
			}
		default:
			e.log.Error(fmt.Sprintf("Unsupported dct format: %s", dist.DctFormat))
		}
	}

	return nil
}
