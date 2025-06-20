package repositories

import (
	"database/sql"
	"errors"
	"vkspam/database"
	"vkspam/models"
)

type DistributionGroupRepository interface {
	GetList(userId int) ([]models.DistributionGroup, error)
	Save(dg models.DistributionGroup) (int, error)
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
	rows, err := d.DB.Query(
		"SELECT "+
			" id,"+
			" name, "+
			" description, "+
			" user_id, "+
			" sex, "+
			" only_birthday_today,"+
			" only_birthday_friends,"+
			" last_processing"+
			" FROM distributiongroup WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var distributionGroups []models.DistributionGroup
	for rows.Next() {
		var distributionGroup models.DistributionGroup
		if err := rows.Scan(
			&distributionGroup.Id,
			&distributionGroup.Name,
			&distributionGroup.Description,
			&distributionGroup.UserId,
			&distributionGroup.Sex,
			&distributionGroup.OnlyBirthdayToday,
			&distributionGroup.OnlyBirthdayFriends,
			&distributionGroup.LastProcessing,
		); err != nil {
			return nil, err
		}

		distributionGroups = append(distributionGroups, distributionGroup)
	}

	return distributionGroups, nil
}

func (d *distributionGroupRepository) Save(dg models.DistributionGroup) (int, error) {
	var id int
	db, _ := database.GetDBInstance()
	if dg.Id > 0 {
		_, err := db.Db.Exec(
			`UPDATE distributiongroup 
				SET name = $1, 
				    description = $2, 
				    sex = $3, 
				    only_birthday_today = $4,
				    only_birthday_friends = $5,
				    last_processing = $6
				WHERE id = $7;`,
			dg.Name, dg.Description, dg.Sex, dg.OnlyBirthdayToday, dg.OnlyBirthdayFriends, dg.LastProcessing, dg.Id)
		if err != nil {
			return 0, errors.New("error insert Distribution group")
		}
		id = dg.Id
	} else {
		err := db.Db.QueryRow(
			`INSERT INTO distributiongroup (
					   name,
					   description, 
					   user_id, 
					   sex, 
					   only_birthday_today, 
					   only_birthday_friends,
					   last_processing
					   ) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			dg.Name, dg.Description, dg.UserId, dg.Sex, dg.OnlyBirthdayToday, dg.OnlyBirthdayFriends, dg.LastProcessing,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
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
	rows, err := d.DB.Query("SELECT"+
		" id,"+
		" name,"+
		" description, "+
		" user_id, "+
		" sex, "+
		" only_birthday_today,"+
		" only_birthday_friends,"+
		" last_processing"+
		" FROM distributiongroup WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var distributionGroup models.DistributionGroup
	for rows.Next() {
		if err := rows.Scan(
			&distributionGroup.Id, &distributionGroup.Name, &distributionGroup.Description,
			&distributionGroup.UserId, &distributionGroup.Sex, &distributionGroup.OnlyBirthdayToday,
			&distributionGroup.OnlyBirthdayFriends, &distributionGroup.LastProcessing,
		); err != nil {
			return nil, err
		}
	}

	return &distributionGroup, nil
}
