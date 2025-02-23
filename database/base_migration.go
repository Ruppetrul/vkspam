package database

import "vkspam/models"

type BaseMigration struct{}

func (v BaseMigration) Run(db models.DbSingleton, sql string) (success bool, err error) {
	_, err = db.Db.Exec(sql)

	if err != nil {
		return false, err
	}

	return true, nil
}
