package models

import "time"

type DistributionGroup struct {
	Id                  int       `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	UserId              int       `json:"user_id"`
	Sex                 int       `json:"sex"`
	OnlyBirthdayToday   bool      `json:"only_birthday_today"`
	OnlyBirthdayFriends bool      `json:"only_birthday_friends"`
	LastProcessing      time.Time `json:"last_processing"`
}
