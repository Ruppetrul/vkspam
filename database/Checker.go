package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"vkspam/database/migrations"
	"vkspam/models"
)

func CheckAndMigrate() {
	db, err := GetDBInstance()
	if err != nil {
		log.Fatal("Db connection error", err)
	}

	runMigration(db, migrations.CreateUsersMigration{}.GetSql())
	runMigration(db, migrations.AddDistributionGroup{}.GetSql())
	runMigration(db, migrations.AddDistribution{}.GetSql())
	runMigration(db, migrations.AddDistributionUrl{}.GetSql())
}

func runMigration(db models.DbSingleton, sql string) {
	for retries := 0; retries < 5; retries++ {
		_, err := BaseMigration{}.Run(db, sql)
		if err == nil {
			return
		}

		if err.Error() == "pq: the database system is starting up" {
			fmt.Println("Database is already running")
			time.Sleep(2 * time.Second)
			continue
		}
		log.Fatalf("Migration error %v", err)
	}
}

var instance models.DbSingleton
var once sync.Once

func GetDBInstance() (models.DbSingleton, error) {
	var err error

	params := map[string]string{
		"user":     os.Getenv("POSTGRES_USER"),
		"password": os.Getenv("POSTGRES_PASSWORD"),
		"dbname":   os.Getenv("POSTGRES_NAME"),
		"sslmode":  os.Getenv("POSTGRES_SSL_MODE"),
		"host":     os.Getenv("POSTGRES_HOST"),
		"port":     os.Getenv("POSTGRES_PORT"),
	}

	var connStr string
	for key, value := range params {
		connStr += fmt.Sprintf("%s=%s ", key, value)
	}
	connStr = strings.TrimSpace(connStr)

	once.Do(func() {
		instance = models.DbSingleton{}
		instance.Db, err = sql.Open("postgres", connStr)

		for retries := 0; retries < 5; retries++ {
			err = instance.Db.Ping()
			if err == nil {
				break
			}

			if err.Error() == "pq: the database system is starting up" {
				fmt.Println("Database is already running")
				time.Sleep(2 * time.Second)
				continue
			}
			log.Fatalf("Error ping database %v", err)
		}
	})

	return instance, nil
}
