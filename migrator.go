package go_migrator

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

// TODO - refactor Up and Down to use the same code
// TODO - implement ErrorInterface for migrator errors, to handle user info in the error message

func (m Migrator) getAppliedMigrationMap() (map[int]bool, error) {
	appliedMigrations, err := m.driver.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	appliedMigrationMap := make(map[int]bool)
	for index := range appliedMigrations {
		version, err := strconv.Atoi(appliedMigrations[index])
		if err != nil {
			return nil, err
		}
		appliedMigrationMap[version] = true
	}

	return appliedMigrationMap, nil
}

func (m Migrator) readMigrationFiles(direction MigrationDirection) ([]string, error) {
	// TODO - update logic to operate on regex based on go_migrator.DefaultFileNameFormat
	expectedFilenameSuffix := fmt.Sprintf(".%s.sql", strings.ToLower(direction.String()))

	files, err := ioutil.ReadDir(m.config.AbsoluteMigrationsDirectory)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for index := range files {
		if files[index].IsDir() {
			continue
		}

		if !strings.HasSuffix(files[index].Name(), expectedFilenameSuffix) {
			continue
		}

		migrationFiles = append(migrationFiles, files[index].Name())
	}

	sort.Strings(migrationFiles)

	return migrationFiles, nil
}

func (m Migrator) getVersionFromFileName(fileName string) (int, error) {
	parts := strings.Split(fileName, m.config.FileNameVersionDelimiter)

	if len(parts) != m.config.PartsCountByDelimiter {
		return 0, fmt.Errorf(`invalid migration file name %s, mismatch part count configured %d vs %d on split by delimiter`, fileName, m.config.PartsCountByDelimiter, len(parts))
	}

	fileNameParts := strings.Split(parts[0], m.config.FileNameVersionDelimiter)

	// TODO - validate timestamp
	version := fileNameParts[m.config.VersionIndexByDelimiter]

	return strconv.Atoi(version)
}

func (m Migrator) Migrate() error {
	err := m.driver.InitSchema()
	if err != nil {
		return err
	}

	appliedMigrationMap, err := m.getAppliedMigrationMap()
	if err != nil {
		return err
	}

	migrationFiles, err2 := m.readMigrationFiles(MigrationDirectionUp)
	if err2 != nil {
		return err2
	}

	for index := range migrationFiles {
		version, err2 := m.getVersionFromFileName(migrationFiles[index])
		if err2 != nil {
			return err2
		}

		if _, ok := appliedMigrationMap[version]; ok {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", m.config.AbsoluteMigrationsDirectory, migrationFiles[index])
		rawContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		err3 := m.driver.Apply(string(rawContent))
		if err3 != nil {
			return err3
		}

		err = m.driver.PushVersion(version)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Migrator) Up() error {
	err := m.driver.InitSchema()
	if err != nil {
		return err
	}

	currentVersion, err := m.driver.GetVersion()
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		currentVersion = 0
	}

	fileName := ""
	fileVersion := 0
	migrationFiles, err2 := m.readMigrationFiles(MigrationDirectionUp)
	if err2 != nil {
		return err2
	}

	for index := range migrationFiles {
		fileVersion, err2 = m.getVersionFromFileName(migrationFiles[index])
		if err2 != nil {
			return err2
		}

		if fileVersion > currentVersion {
			fileName = migrationFiles[index]
			break
		}
	}

	if fileName == "" {
		return fmt.Errorf("no un-applied migration files found after the current version %d", currentVersion)
	}

	filePath := fmt.Sprintf("%s/%s", m.config.AbsoluteMigrationsDirectory, fileName)
	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err3 := m.driver.Apply(string(rawContent))
	if err3 != nil {
		return err3
	}

	err = m.driver.PushVersion(fileVersion)
	if err != nil {
		return err
	}

	return nil
}

func (m Migrator) Down() error {
	err := m.driver.InitSchema()
	if err != nil {
		return err
	}

	currentVersion, err := m.driver.GetVersion()

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("current version not found to migrate down, please check if there are any applied migrations")
		}

		return err
	}

	fileName := ""
	fileVersion := 0
	migrationFiles, err2 := m.readMigrationFiles(MigrationDirectionDown)
	if err2 != nil {
		return err2
	}

	for index := range migrationFiles {
		fileVersion, err2 = m.getVersionFromFileName(migrationFiles[index])
		if err2 != nil {
			return err2
		}

		if fileVersion == currentVersion {
			fileName = migrationFiles[index]
			break
		}
	}

	if fileName == "" {
		return fmt.Errorf("no migration files found for the current version %d", currentVersion)
	}

	filePath := fmt.Sprintf("%s/%s", m.config.AbsoluteMigrationsDirectory, fileName)
	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err3 := m.driver.Apply(string(rawContent))
	if err3 != nil {
		return err3
	}

	err = m.driver.PopVersion()
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) Version() (int, error) {
	err := m.driver.InitSchema()
	if err != nil {
		return 0, err
	}

	version, err := m.driver.GetVersion()
	if err != nil {
		panic(err)
	}

	return version, nil
}

func (m Migrator) Execute(version int, direction MigrationDirection) error {
	fileName := ""
	fileVersion := 0
	migrationFiles, err := m.readMigrationFiles(direction)
	if err != nil {
		return err
	}

	for index := range migrationFiles {
		if index == version {
			fileName = migrationFiles[index]
			break
		}
	}
	for index := range migrationFiles {
		fileVersion, err = m.getVersionFromFileName(migrationFiles[index])
		if err != nil {
			return err
		}

		if fileVersion == version {
			fileName = migrationFiles[index]
			break
		}
	}

	if fileName == "" {
		return fmt.Errorf("no migration files found for the current version %d for %s", version, direction)
	}

	filePath := fmt.Sprintf("%s/%s", m.config.AbsoluteMigrationsDirectory, fileName)
	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = m.driver.Apply(string(rawContent))
	if err != nil {
		return err
	}

	return nil
}
