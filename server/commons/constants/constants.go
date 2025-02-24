// Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.

package constants

const (
	ComponentName                = "meeting-analyzer"
	DatabaseName                 = "meeting_analyzer_db"
	ContextTimeout               = 5
	HTTPServerRequestReadTimeOut = 5
	DefaultPort                  = "8080"
	FileNotation                 = "file://"
	MigrationFolderPath          = FileNotation + "server/repositories/db/migrations"
	EnvVarDBPA                   = "POSTGRES_PASSWORD"
	EnvVarDBUser                 = "POSTGRES_USER"
	EnvVarDBHost                 = "POSTGRES_HOST"
	EnvVarDBPort                 = "POSTGRES_PORT"
)
