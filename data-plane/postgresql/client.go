package postgresql

import (
	"fmt"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/data-plane"
	"github.com/YasiruR/go-connector/domain/models"
	"github.com/YasiruR/go-connector/domain/pkg"
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

type Client struct {
	host string
	db   string
	usr  string
	pw   string
	log  pkg.Log
}

func NewClient(cfg boot.Config, log pkg.Log) *Client {
	var e Client
	for _, ds := range cfg.Databases {
		if ds.Name == data_plane.DatabasePostgresql {
			e = Client{
				host: ds.Hostname,
				db:   ds.Database,
				usr:  ds.Username,
				pw:   ds.Password,
				log:  log,
			}
		}
	}

	return &e
}

func (d *Client) Dump(dest models.Database) error {
	fileName := fmt.Sprintf(data_plane.PushFilePrefix+`%s_%d.sql`, dest.Endpoint, time.Now().Unix())

	out, err := exec.Command(`bash`, `-c`, fmt.Sprintf(
		dumpCmd, d.pw, d.host, d.usr, d.db, fileName)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run dump command: %w (stderr: %s)", err, string(out))
	}

	out, err = exec.Command(`bash`, `-c`, fmt.Sprintf(
		restoreCmd, dest.Password, dest.Endpoint, dest.User, dest.Name, fileName)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run restore command: %w (stderr: %s)", err, string(out))
	}

	d.log.Info(fmt.Sprintf("dataset migrated to postgresql instance (endpoint: %s, database: %s)",
		dest.Endpoint, dest.Name))
	return nil
}

func (d *Client) Store(src models.Database) error {
	fileName := fmt.Sprintf(data_plane.PullFilePrefix+`%s_%d.sql`, src.Endpoint, time.Now().Unix())

	out, err := exec.Command(`bash`, `-c`, fmt.Sprintf(
		dumpCmd, src.Password, src.Endpoint, src.User, src.Name, fileName)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run dump command: %w (stderr: %s)", err, string(out))
	}

	out, err = exec.Command(`bash`, `-c`, fmt.Sprintf(
		restoreCmd, d.pw, d.host, d.usr, d.db, fileName)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run restore command: %w (stderr: %s)", err, string(out))
	}

	d.log.Info(fmt.Sprintf("dataset fetched from postgresql instance (endpoint: %s, database: %s)",
		src.Endpoint, src.Name))
	return nil
}
