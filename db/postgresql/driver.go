package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" // postgresql driver implicitly required
)

func (p DriverConfig) GetDriverName() string {
	return "postgres"
}

func (p DriverConfig) GetDSN() string {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		p.GetDriverName(),
		p.Username,
		p.Password,
		p.Host,
		p.Port,
		p.DatabaseName,
	)

	if p.SSL {
		dsn += fmt.Sprintf("?sslmode=require")
	} else {
		dsn += fmt.Sprintf("?sslmode=disable")
	}

	return dsn
}

func (p Driver) Connect() error {
	var err error
	p.client, err = sql.Open(p.config.GetDriverName(), p.config.GetDSN())
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
		(version BIGINT NOT NULL UNIQUE, applied_at TIMESTAMP WITH TIME ZONE NOT NULL)`,
		p.config.MigrationTableName)

	_, err := p.client.Exec(query)

	return err
}

func (p Driver) GetVersion() (int, error) {
	var version int

	query := fmt.Sprintf(`SELECT version FROM %s ORDER BY applied_at DESC LIMIT 1`, p.config.MigrationTableName)
	err := p.client.QueryRow(query).Scan(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func (p Driver) PushVersion(version int) error {
	query := fmt.Sprintf(`INSERT INTO  %s (version, applied_at) VALUES($1, NOW())`, p.config.MigrationTableName)

	_, err := p.client.Exec(query, version)

	return err
}

func (p Driver) PopVersion() error {
	query := `DELETE FROM migrations WHERE version = (SELECT version FROM migrations ORDER BY applied_at DESC LIMIT 1)`

	_, err := p.client.Exec(query)

	return err
}

func (p Driver) GetAppliedMigrations() ([]string, error) {
	var versions []string

	query := fmt.Sprintf(`SELECT version FROM %s ORDER BY applied_at DESC`, p.config.MigrationTableName)

	rows, err := p.client.Query(query)
	if err != nil {
		return versions, err
	}

	var version string
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			return versions, err
		}

		versions = append(versions, version)
	}

	if rows.Err() != nil {
		return versions, rows.Err()
	}

	return versions, nil
}

func PrettyPrint(prefix string, i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(prefix, string(s))
}
