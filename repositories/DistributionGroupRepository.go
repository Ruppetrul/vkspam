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
	Delete(id int) error
	Get(id int) (*models.DistributionGroup, error)
}

type distributionGroupRepository struct {
	DB *sql.DB
}

func NewDistributionGroupRepository(db *sql.DB) DistributionGroupRepository {
	return &distributionGroupRepository{DB: db}
}

func (d *distributionGroupRepository) GetList(userId int) ([]models.DistributionGroup, error) {
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

func (d *distributionGroupRepository) Save(dg models.DistributionGroup) error {
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

func (d *distributionGroupRepository) Delete(id int) error {
	result, err := d.DB.Exec("DELETE FROM distributiongroup WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no record found to delete")
	}

	return nil
}

func (d *distributionGroupRepository) Get(id int) (*models.DistributionGroup, error) {
	rows, err := d.DB.Query("SELECT * FROM distributiongroup WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var distributionGroup models.DistributionGroup
	for rows.Next() {
		if err := rows.Scan(&distributionGroup.Id, &distributionGroup.Name, &distributionGroup.Description, &distributionGroup.UserId); err != nil {
			return nil, err
		}
	}

	return &distributionGroup, nil
}
