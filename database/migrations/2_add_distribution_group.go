package migrations

type AddDistributionGroup struct {
}

func (v AddDistributionGroup) GetSql() (sql string) {
	return `
		CREATE TABLE IF NOT EXISTS DistributionGroup (
			id SERIAL PRIMARY KEY,
			name VARCHAR(256) NOT NULL,
			description TEXT,
		    user_id INT
		);
	`
}
