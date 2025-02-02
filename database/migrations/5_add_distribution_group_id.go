package migrations

type AddDistributionNumber struct {
}

func (v AddDistributionNumber) GetSql() (sql string) {
	return `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                           WHERE table_name='distribution' AND column_name='group_id') THEN
				ALTER TABLE distribution ADD COLUMN group_id integer;
			END IF;
		END $$;
	`
}
