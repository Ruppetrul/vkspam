package migrations

type AddVkToken struct {
}

func (v AddVkToken) GetSql() (sql string) {
	return `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                           WHERE table_name='users' AND column_name='vk_token') THEN
				ALTER TABLE users ADD COLUMN vk_token text;
			END IF;
		END $$;
	`
}
