package models

type User struct {
	Id         int    `db:"id"`
	Email      string `db:"email"`
	InviteCode string `db:"invite_code"`
	Password   string `db:"password"`
	Token      string `db:"token"`
}
