package postgresql

import (
	"fmt"
	"github.com/YasiruR/go-connector/domain/boot"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os/exec"
)

// pg_dump -h <src-host> -U <src-user> <src-database> > backup.sql
// psql -h <dest-host> -U <dest-user> <dest-database> < backup.sql

const (
	dumpCmd    = "PGPASSWORD=%s pg_dump -h %s -U %s %s > backup.sql"
	restoreCmd = "PGPASSWORD=%s psql -h %s -U %s %s < backup.sql"
)

type Exchanger struct {
	host string
	db   string
	usr  string
	pw   string
}

func NewExchanger(cfg boot.Config) *Exchanger {
	var e Exchanger
	for _, ds := range cfg.DataSources {
		if ds.Name == `postgresql` {
			e = Exchanger{
				host: ds.Hostname,
				db:   ds.Database,
				usr:  ds.Username,
				pw:   ds.Password,
			}
		}
	}

	return &e
}

func (d *Exchanger) Migrate(host, db, usr, pw string) error {
	if _, err := exec.Command(fmt.Sprintf(dumpCmd, d.pw, host, d.usr, db)).Output(); err != nil {
		return fmt.Errorf("failed to run dump command: %w", err)
	}

	if _, err := exec.Command(fmt.Sprintf(restoreCmd, pw, host, usr, db)).Output(); err != nil {
		return fmt.Errorf("failed to run restore command: %w", err)
	}

	return nil
}
