package repositories

import (
	"database/sql"
	"vkspam/models"
)

type DistributionGroupRepository interface {
	GetList(userId int) ([]models.DistributionGroup, error)
}

type distributionRepository struct {
	DB *sql.DB
}

func NewDistributionRepository(db *sql.DB) DistributionGroupRepository {
	return &distributionRepository{DB: db}
}

func (d *distributionRepository) GetList(userId int) ([]models.DistributionGroup, error) {
	rows, err := d.DB.Query("SELECT * FROM distributiongroup WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var distributionGroups []models.DistributionGroup
	for rows.Next() {
		var distributionGroup models.DistributionGroup
		if err := rows.Scan(&distributionGroup.Id, &distributionGroup.Name, &distributionGroup.Description, &distributionGroup.UserId); err != nil {
			return nil, err
		}

		distributionGroups = append(distributionGroups, distributionGroup)
	}

	return distributionGroups, nil
}
