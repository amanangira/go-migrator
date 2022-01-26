package postgresql

import (
	"database/sql"
	"fmt"
)

func (p DriverConfig) GetDriverName() string {
	return "postgres"
}

func (p DriverConfig) GetDSN() string {
	var sslMode string
	if p.SSL {
		sslMode = "sslmode=require"
	} else {
		sslMode = "sslmode=disable"
	}

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		p.GetDriverName(),
		p.Username,
		p.Password,
		p.Host,
		p.Port,
		p.DatabaseName,
		sslMode,
	)

	return dsn
}

func (p Driver) Connect() error {
	var err error
	p.client, err = sql.Open("postgres", p.config.GetDSN())
	if err != nil {
		return err
	}

	return nil
}

func (p Driver) Close() error {
	return p.client.Close()
}

func (p Driver) GetClient() *sql.DB {
	return p.client
}

func (p Driver) Apply(command string) error {
	_, err := p.client.Exec(command)

	return err
}

func (p Driver) InitSchema() error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
		(version TEXT NOT NULL, applied_at TIMESTAMP WITH TIME ZONE NOT NULL)`,
		p.config.MigrationTableName)

	_, err := p.client.Exec(query)

	return err
}

func (p Driver) GetVersion() (string, error) {
	var version string

	query := fmt.Sprintf(`SELECT version FROM %s ORDER BY applied_at DESC LIMIT 1`, p.config.MigrationTableName)
	err := p.client.QueryRow(query).Scan(&version)
	if err != nil {
		return "", err
	}

	return version, nil
}

func (p Driver) PushVersion(version string) error {
	query := fmt.Sprintf(`INSERT INTO  %s (version, applied_at) VALUES($1, NOW())`, p.config.MigrationTableName)

	_, err := p.client.Exec(query, version)

	return err
}

func (p Driver) PopVersion(version string) error {
	query := `DELETE FROM migrations WHERE version=$1`

	_, err := p.client.Exec(query, version)

	return err
}
