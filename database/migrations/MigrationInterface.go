package migrations

import "vkspam/database"

type MigrationInterface interface {
	Run(db database.DbSingleton) (success bool, err error)
}
