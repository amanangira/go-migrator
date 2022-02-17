package postgresql

import (
	"database/sql"
)

type DriverConfig struct {
	Host               string
	Port               int
	Username           string
	SSL                bool
	Password           string
	DatabaseName       string
	MigrationTableName string
}

type Driver struct {
	config DriverConfig
	client *sql.DB
}

func NewConfig(
	host string,
	port int,
	user string,
	ssl bool,
	password,
	database,
	migrationTableName string,
) DriverConfig {
	return DriverConfig{
		Host:         host,
		Port:         port,
		Username:     user,
		SSL:          ssl,
		Password:     password,
		DatabaseName: database,

		// TODO - allow overriding of this value via interface implementation
		MigrationTableName: migrationTableName,
	}
}

func NewDriverWithConfig(config DriverConfig) (Driver, error) {
	driver := Driver{}

	dbClient, err := sql.Open(config.GetDriverName(), config.GetDSN())
	if err != nil {
		return driver, err
	}

	driver = Driver{
		config: config,
		client: dbClient,
	}

	return driver, nil
}

func NewDriver(
	host string,
	port int,
	user string,
	ssl bool,
	password,
	database,
	migrationTableName string,
) (Driver, error) {
	config := DriverConfig{
		Host:         host,
		Port:         port,
		Username:     user,
		SSL:          ssl,
		Password:     password,
		DatabaseName: database,

		// TODO - allow overriding of this value via interface implementation
		MigrationTableName: migrationTableName,
	}

	return NewDriverWithConfig(config)
}
