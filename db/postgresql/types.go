package postgresql

import "database/sql"

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
) *DriverConfig {
	return &DriverConfig{
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

func NewDriver(config *DriverConfig) *Driver {
	return &Driver{
		config: *config,
	}
}
