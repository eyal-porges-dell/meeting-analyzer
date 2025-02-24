// Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.

package db

import (
	"database/sql"
	"meeting-analyzer/server/repositories"

	"eos2git.cec.lab.emc.com/ISG-Edge/hzp-go-commons/database/interfaces"
)

type repository struct {
	db    interfaces.Database
	dbCon *sql.DB
}

// NewRepository creates a new DriftRepository with dependency injection
func NewRepository(con interfaces.Database) (repositories.Repository, error) {
	//dbCon, err := con.GetConnection()
	//if err != nil {
	//	return nil, err
	//}

	return &repository{
		db:    nil,
		dbCon: nil,
	}, nil
}
