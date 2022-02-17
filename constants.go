package go_migrator

const (
	DefaultPartsCountForFileName = 3
	DirectionTextUp              = "up"
	DirectionTextDown            = "down"
	MigrationFileNameHelpText    = "%s invalid migration file name format, should be of format {version}_{title}.{up|down}.{extension}"

	MigrationDirectionUp MigrationDirection = iota
	MigrationDirectionDown

	FileNameDefaultExtensionDelimiter = "."
	FileNameVersionPlaceholder        = "{version}"
	FileNameDirectionPlaceholder      = "{direction}"
	FileNameExtensionPlaceholder      = "{extension}"

	CommandUsageHelpText = "Description:\n  Execute a migration to a specified version or the latest available version." +
		"\n\nUsage:\n  ./main [options] [<version>] [<direction>] " +
		"\n\nArguments:\n  version                     The version number (YYYYMMDDHHMMSS) to migrate to." +
		"\n  direction                   The direction to apply a migration towards. Possible values - up, down" +
		"\n\nOptions:\n      migrate                Applies all the unapplied migrations in the ascending order." +
		"\n      up                     Rolls up the migration by one version." +
		"\n      down                   Rolls down the migration by one version." +
		"\n      execute                Applies the specified migration. No warning is given even if the migration is already applied. " +
		"Usage - ./main execute <version> <direction>" +
		"\n      version                Displays the current DB migration version."
)
