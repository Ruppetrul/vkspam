package migrations

type CreateUsersMigration struct {
}

func (v CreateUsersMigration) GetSql() (sql string) {
	return `
        CREATE TABLE IF NOT EXISTS users (
            id serial NOT NULL,
            email VARCHAR(256) NOT NULL,
            password VARCHAR(256) NOT NULL,
            PRIMARY KEY (id)
        );

        INSERT INTO users (email, password)
        SELECT 'admin@gmail.com', '$2a$10$q1P4dpQYJO/kgbD6V0/bU.wIEpFvhQZ6sETxaShy9dA8UUIZ6z8ZS'
        WHERE NOT EXISTS (
            SELECT 1 FROM users WHERE email = 'admin@gmail.com'
        );
    `
}
