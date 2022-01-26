package go_migrator

const (
	DefaultPartsCountForFileName = 3
	DirectionTextUp              = "up"
	DirectionTextDown            = "down"
	MigrationFileNameHelpText    = "%s invalid migration file name format, should be of format {version}_{title}.{up|down}.{extension}"

	MigrationDirectionUp MigrationDirection = iota
	MigrationDirectionDown
)
