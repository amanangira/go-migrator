package db

import "database/sql"

type MigrationTableName interface {
	GetMigrationTableName() string
}

type DriverInterface interface {
	Connect() error
	Close() error
	GetClient() *sql.DB
	Apply(command string) error
	InitSchema() error
	GetVersion() (int, error)
	GetAppliedMigrations() ([]string, error)
	PushVersion(version int) error
	PopVersion() error
}
