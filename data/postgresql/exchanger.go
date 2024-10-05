package postgresql

import (
	"fmt"
	"github.com/YasiruR/go-connector/domain/boot"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os/exec"
	"time"
)

// pg_dump -c -h <src-host> -U <src-user> <src-database> > backup.sql
// psql -h <dest-host> -U <dest-user> <dest-database> < backup.sql

const (
	dumpCmd    = "PGPASSWORD=%s pg_dump -c -h %s -U %s %s > %s"
	restoreCmd = "PGPASSWORD=%s psql -h %s -U %s %s < %s"
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
	fileName := fmt.Sprintf(`backup_%s_%d.sql`, host, time.Now().Unix())
	cmd := fmt.Sprintf(dumpCmd, d.pw, d.host, d.usr, d.db, fileName)
	out, err := exec.Command(`bash`, `-c`, cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run dump command: %w (stderr: %s)", err, string(out))
	}

	cmd = fmt.Sprintf(restoreCmd, pw, host, usr, db, fileName)
	out, err = exec.Command(`bash`, `-c`, cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run restore command: %w (stderr: %s)", err, string(out))
	}

	return nil
}
