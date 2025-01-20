package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"vkspam/database/migrations"
)

func CheckAndMigrate() {
	db, err := GetDBInstance()
	if err != nil {
		log.Fatal("Db connection error", err)
	}

	runMigration(db, migrations.CreateUsersMigration{})
}

func runMigration(db DbSingleton, migration migrations.MigrationInterface) {
	_, err := migration.Run(db)
	if err != nil {
		log.Fatal("Migration error", err)
	}
}

type DbSingleton struct {
	Db *sql.DB
}

var instance DbSingleton
var once sync.Once

func GetDBInstance() (DbSingleton, error) {
	var err error

	params := map[string]string{
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbname":   os.Getenv("DB_NAME"),
		"sslmode":  os.Getenv("DB_SSL_MODE"),
		"host":     os.Getenv("DB_HOST"),
	}

	var connStr string
	for key, value := range params {
		connStr += fmt.Sprintf("%s=%s ", key, value)
	}
	connStr = strings.TrimSpace(connStr)

	once.Do(func() {
		instance = DbSingleton{}
		instance.Db, err = sql.Open("postgres", connStr)
	})

	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	return instance, nil
}
