package migrations

import "vkspam/database"

type CreateUsersMigration struct{}

func (v CreateUsersMigration) Run(db database.DbSingleton) (success bool, err error) {
	_, err = db.Db.Exec(`CREATE TABLE IF NOT EXISTS users (id integer NOT NULL, PRIMARY KEY (id));`)

	if err != nil {
		return false, err
	}

	return true, nil
}
