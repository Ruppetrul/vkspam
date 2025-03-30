package migrations

type AddLastProcessing struct {
}

func (v AddLastProcessing) GetSql() (sql string) {
	return `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                           WHERE table_name='distributiongroup' AND column_name='last_processing') THEN
				ALTER TABLE distributiongroup ADD COLUMN last_processing timestamp NULL;
			END IF;
		END $$;
	`
}
