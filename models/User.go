package models

type User struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	InviteCode string `db:"invite_code"`
	Token      string `db:"token"`
}
