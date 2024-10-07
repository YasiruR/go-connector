package data_plane

import (
	"fmt"
	"github.com/YasiruR/go-connector/data-plane/postgresql"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/data-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
	"os/exec"
)

// Exchanger can invoke external service clients (containers) for data-plane transfer
type Exchanger struct {
	psql *postgresql.Client
	cat  stores.ProviderCatalog
}

func NewExchanger(cfg boot.Config, s domain.Stores, log pkg.Log) *Exchanger {
	// create new directory for data-plane backups
	if _, err := exec.Command(`bash`, `-c`, `mkdir -p backups/pull && mkdir -p backups/push`).
		CombinedOutput(); err != nil {
		log.Fatal(`creating backup folders failed`, err)
	}

	return &Exchanger{
		cat:  s.ProviderCatalog,
		psql: postgresql.NewClient(cfg, log),
	}
}

func (e *Exchanger) NewToken(participantId, datasetId string) string {
	return ``
}

func (e *Exchanger) PushWithCredentials(et data_plane.ExchangerType, dest data_plane.Database) error {
	switch et {
	case data_plane.TypePostgresql:
		if err := e.psql.Dump(dest); err != nil {
			return fmt.Errorf("postgresql failed - %s", err)
		}
	default:
		return errors.Client(errors.ExchangerError(`no supported data-plane distribution type found`))
	}

	return nil
}

func (e *Exchanger) PullWithCredentials(src data_plane.Database) error {
	return nil
}
