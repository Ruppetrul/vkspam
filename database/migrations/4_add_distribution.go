package migrations

import "vkspam/models"

type AddDistribution struct {
}

func (v AddDistribution) GetSql() (sql string) {
	var anyPublic models.DistributionType = models.AnyPublic

	return `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'distribution_type') THEN
				CREATE TYPE distribution_type AS ENUM ('` + anyPublic.String() + `');
			END IF;
		END $$;
		CREATE TABLE IF NOT EXISTS distribution (
			id SERIAL PRIMARY KEY,
			name VARCHAR(256) NOT NULL,
			type distribution_type
		);
	`
}
