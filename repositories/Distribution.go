package repositories

import (
	"database/sql"
	"errors"
	"vkspam/database"
	"vkspam/models"
)

type DistributionRepository interface {
	Save(d models.Distribution) error
	Get(id int) (*models.Distribution, error)
}

type distributionRepository struct {
	DB *sql.DB
}

func NewDistributionRepository(db *sql.DB) DistributionRepository {
	return &distributionRepository{DB: db}
}

func (dr *distributionRepository) Save(d models.Distribution) error {
	if d.Id > 0 {
		db, _ := database.GetDBInstance()
		_, err := db.Db.Exec(`UPDATE distribution SET name = $1, type = $2, url = $3 WHERE id = $4;`, d.Name, d.Type, d.Url)
		if err != nil {
			return errors.New("Error insert Distribution")
		}
	} else {
		db, _ := database.GetDBInstance()
		err := db.Db.QueryRow(`INSERT INTO distribution (name, type, url) VALUES ($1, $2, $3)`, d.Name, d.Type, d.Url)
		if err != nil {
			return err.Err()
		}
	}

	return nil
}

func (d *distributionRepository) Get(id int) (*models.Distribution, error) {
	rows, err := d.DB.Query("SELECT * FROM distribution WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var distribution models.Distribution
	for rows.Next() {
		if err := rows.Scan(&distribution.Id, &distribution.Name, &distribution.Type, &distribution.Url, &distribution.GroupId); err != nil {
			return nil, err
		}
	}

	return &distribution, nil
}
