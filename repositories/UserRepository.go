package repositories

import (
	"database/sql"
	"errors"
	"vkspam/database"
	"vkspam/models"
)

type UserRepository interface {
	TryLogin(email string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	Save(user *models.User) error
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (repo *userRepository) TryLogin(email string) (*models.User, error) {
	rows, err := repo.DB.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user models.User

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (repo *userRepository) FindUserByEmail(email string) (*models.User, error) {
	rows, err := repo.DB.Query("SELECT * FROM users WHERE email = $1;", email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user models.User

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (repo *userRepository) Save(user *models.User) error {
	if user.Id > 0 {
		db, _ := database.GetDBInstance()
		_, err := db.Db.Exec(`UPDATE users SET email = $1, password = $2 WHERE id = $3;`, user.Email, user.Password, user.Id)
		if err != nil {
			return errors.New("error when update user")
		}
	} else {
		db, _ := database.GetDBInstance()
		err := db.Db.QueryRow(`INSERT INTO users (email, password) VALUES ($1, $2)`,
			user.Email, user.Password)
		if err != nil {
			return err.Err()
		}
	}

	return nil
}
