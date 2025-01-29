package repositories

import (
	"database/sql"
)

type DistributionRepository interface {
}

type distributionRepository struct {
	DB *sql.DB
}

func NewDistributionRepository(db *sql.DB) DistributionRepository {
	return &distributionRepository{DB: db}
}
