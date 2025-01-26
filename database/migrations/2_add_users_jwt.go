package migrations

type CreateUsersJwt struct {
}

func (v CreateUsersJwt) GetSql() (sql string) {
	return `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                           WHERE table_name='users' AND column_name='jwt_token') THEN
				ALTER TABLE users ADD COLUMN jwt_token text;
			END IF;
		END $$;
	`
}
