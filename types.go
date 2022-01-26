package go_migrator

import (
	"fmt"
	"path/filepath"

	"github.com/amanangira/go-migrator/db"
)

type MigrationDirection int
type MigrationCommand string

type MigrationInstruction struct {
	Version   int
	Direction MigrationDirection
	Command   MigrationCommand
}

type MigratorConfig struct {
	AbsoluteMigrationsDirectory string

	// Filename config
	// TODO - allow overriding of this value via interface implementation
	FileNameDelimiter       string
	VersionIndexByDelimiter int
	PartsCountByDelimiter   int
}

type Migrator struct {
	config *MigratorConfig
	driver *db.DriverInterface
}

func (m MigrationDirection) String() string {
	switch m {
	case MigrationDirectionUp:
		return "Up"

	case MigrationDirectionDown:
		return "Down"
	}

	panic("unsupported direction")
}

func NewConfig(
	directory string,
) *MigratorConfig {
	sanitizedDirectory := sanitizeDirectoryPath(directory)

	return &MigratorConfig{
		AbsoluteMigrationsDirectory: sanitizedDirectory,
		FileNameDelimiter:           ".",
		VersionIndexByDelimiter:     0,
		PartsCountByDelimiter:       DefaultPartsCountForFileName,
	}
}

func NewMigrator(
	config *MigratorConfig,
	driver *db.DriverInterface,
) *Migrator {
	return &Migrator{
		config: config,
		driver: driver,
	}
}

func sanitizeDirectoryPath(directory string) string {
	if directory == "" {
		panic(fmt.Sprintf(`invalid directory path %s`, directory))
	}

	directory = filepath.Clean(directory)

	if directory[0:1] == "/" {
		return directory
	}

	var err error
	directory, err = filepath.Abs(directory)

	if err != nil {
		panic(fmt.Sprintf(`error generating absolute path for directory %s and error %s`, directory, err.Error()))
	}

	return directory
}
