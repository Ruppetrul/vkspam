package migrations

type CreateUsersMigration struct {
}

func (v CreateUsersMigration) GetSql() (sql string) {
	return "CREATE TABLE IF NOT EXISTS users (" +
		" id serial NOT NULL," +
		" email VARCHAR(256) NOT NULL," +
		" password VARCHAR(256) NOT NULL," +
		" PRIMARY KEY (id));"
}
