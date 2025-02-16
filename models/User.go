package models

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	VkToken  string `db:"vk_token"`
}
