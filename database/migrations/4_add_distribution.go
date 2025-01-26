package migrations

import "vkspam/models"

type AddDistribution struct {
}

func (v AddDistribution) GetSql() (sql string) {
	var anyPublic models.DistributionType = models.AnyPublic

	return `
		CREATE TYPE distribution_type AS ENUM ('` + anyPublic.String() + `');
		CREATE TABLE IF NOT EXISTS distribution (
			id SERIAL PRIMARY KEY,
			name VARCHAR(256) NOT NULL,
			type distribution_type
		);
	`
}
