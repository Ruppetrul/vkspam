package database

import "vkspam/models"

type MigrationInterface interface {
	Run(db models.DbSingleton) (success bool, err error)
}
