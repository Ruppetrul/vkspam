package database

type MigrationInterface interface {
	GetSql() string
}
