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
	GetVersion() (string, error)
	PushVersion(version string) error
	PopVersion(version string) error
}
