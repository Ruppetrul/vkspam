package models

import "database/sql"

type DbSingleton struct {
	Db *sql.DB
}
