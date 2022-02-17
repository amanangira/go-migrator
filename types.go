package go_migrator

import (
	"fmt"
	"path/filepath"

	"github.com/amanangira/go-migrator/db"
)

type MigrationDirection int

type MigrationInstruction struct {
	Version   int
	Direction MigrationDirection
	Command   string
}

type Config struct {
	AbsoluteMigrationsDirectory string
	Debug                       bool

	// Filename config
	// TODO - allow overriding of this value via interface implementation
	// TODO - Filename config made part of this config in order to allow customization in future
	FileNameVersionDelimiter string
	VersionIndexByDelimiter  int
	PartsCountByDelimiter    int
	FileNameFormat           string
}

type Migrator struct {
	config Config
	driver db.DriverInterface
}

func NewConfig(
	directory string,
	debug bool,
) Config {
	sanitizedDirectory := sanitizeDirectoryPath(directory)

	return Config{
		Debug:                       debug,
		AbsoluteMigrationsDirectory: sanitizedDirectory,
		FileNameVersionDelimiter:    FileNameDefaultExtensionDelimiter,
		VersionIndexByDelimiter:     0,
		PartsCountByDelimiter:       DefaultPartsCountForFileName,
		FileNameFormat:              DefaultFileNameFormat(),
	}
}

func DefaultFileNameFormat() string {
	return fmt.Sprintf(`%s%s%s%s%s`,
		FileNameVersionPlaceholder,
		FileNameDefaultExtensionDelimiter,
		FileNameDirectionPlaceholder,
		FileNameDefaultExtensionDelimiter,
		FileNameExtensionPlaceholder)
}

func NewMigratorWithConfig(
	config Config,
	driver db.DriverInterface,
) *Migrator {
	return &Migrator{
		config: config,
		driver: driver,
	}
}

func NewMigrator(
	directory string,
	debug bool,
	driver db.DriverInterface,
) (*Migrator, error) {
	sanitizedDirectory := sanitizeDirectoryPath(directory)

	config := Config{
		Debug:                       debug,
		AbsoluteMigrationsDirectory: sanitizedDirectory,
		FileNameVersionDelimiter:    FileNameDefaultExtensionDelimiter,
		VersionIndexByDelimiter:     0,
		PartsCountByDelimiter:       DefaultPartsCountForFileName,
		FileNameFormat:              DefaultFileNameFormat(),
	}

	migrator := NewMigratorWithConfig(config, driver)

	err := migrator.driver.Connect()

	return migrator, err
}

func (m MigrationDirection) String() string {
	switch m {
	case MigrationDirectionUp:
		return "Up"

	case MigrationDirectionDown:
		return "Down"
	}

	return ""
}

func (m MigrationDirection) IsValid() bool {
	switch m {
	case MigrationDirectionUp:
		return true

	case MigrationDirectionDown:
		return true
	}

	return false
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
