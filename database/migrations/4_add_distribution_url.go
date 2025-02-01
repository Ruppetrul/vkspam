package migrations

type AddDistributionUrl struct {
}

func (v AddDistributionUrl) GetSql() (sql string) {
	return `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                           WHERE table_name='distribution' AND column_name='url') THEN
				ALTER TABLE distribution ADD COLUMN url text;
			END IF;
		END $$;
	`
}
