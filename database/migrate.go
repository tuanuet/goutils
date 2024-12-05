package migrate

import (
	"fmt"
	"github.com/tuanuet/goutils/logger"
	"os"
	"strconv"
	"strings"
	"time"

	migrateV4 "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const versionTimeFormat = "20060102150405"

// New with up and down options.
func New(migrationFolder string, name string) {
	folder := strings.ReplaceAll(migrationFolder, "file://", "")
	now := time.Now()
	ver := now.Format(versionTimeFormat)

	up := fmt.Sprintf("%s/%s_%s.up.sql", folder, ver, name)
	down := fmt.Sprintf("%s/%s_%s.down.sql", folder, ver, name)

	logger.Infof("create migrate: %s", name)

	if err := os.WriteFile(up, []byte{}, 0644); err != nil {
		logger.Fatalf("create migrate up error: %v", err)
	}
	if err := os.WriteFile(down, []byte{}, 0644); err != nil {
		logger.Fatalf("create migrate down error: %v", err)
	}
}

// Up migrate db to latest version.
func Up(databaseURL string, migrationFolder string) {
	m, err := migrateV4.New(migrationFolder, databaseURL)
	if err != nil {
		logger.Fatalf("error create migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrateV4.ErrNoChange {
		logger.Fatalf("error when migrate up: %v", err)
	}

	logger.Info("migrate up completed")
}

// Down ...
func Down(databaseURL string, migrationFolder string, strVersion string) {
	m, err := migrateV4.New(migrationFolder, databaseURL)

	if err != nil {
		logger.Fatalf("error create migrate: %v", err)
	}

	version, err := strconv.Atoi(strVersion)
	if err != nil {
		logger.Fatalf("error when migrate down: %v", err)
	}

	logger.Infof("migrate down %d", version)
	if err := m.Steps(-version); err != nil {
		logger.Fatalf("error when migrate down: %v", err.Error())
	}
}

// Force ...
func Force(databaseURL string, migrationFolder string, strVersion string) {
	m, err := migrateV4.New(migrationFolder, databaseURL)

	if err != nil {
		logger.Fatalf("error create migrate: %v", err)
	}

	version, err := strconv.Atoi(strVersion)
	if err != nil {
		logger.Fatalf("error when force migrate: %v", err)
	}

	logger.Infof("force to version: %d", version)

	if err := m.Force(version); err != nil {
		logger.Fatalf("error when force db to version %d: %v", version, err)
	}
}
