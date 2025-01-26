package migrations

type CreateUsersMigration struct {
}

func (v CreateUsersMigration) GetSql() (sql string) {
	return "CREATE TABLE IF NOT EXISTS users (id integer NOT NULL, PRIMARY KEY (id));"
}
