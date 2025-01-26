package repositories

import (
	"database/sql"
	"errors"
	"vkspam/database"
	"vkspam/models"
)

type DistributionGroupRepository interface {
	GetList(userId int) ([]models.DistributionGroup, error)
	Save(dg models.DistributionGroup) error
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

func (d *distributionRepository) Save(dg models.DistributionGroup) error {
	if dg.Id > 0 {
		db, _ := database.GetDBInstance()
		_, err := db.Db.Exec(`UPDATE distributiongroup SET name = $1, description = $2 WHERE id = $3;`, dg.Name, dg.Description, dg.Id)
		if err != nil {
			return errors.New("Error insert Distribution group")
		}
	} else {
		db, _ := database.GetDBInstance()
		err := db.Db.QueryRow(`INSERT INTO distributiongroup (name, description, user_id) VALUES ($1, $2, $3)`, dg.Name, dg.Description, dg.UserId)
		if err != nil {
			return err.Err()
		}
	}

	return nil
}
