package models

type Distribution struct {
	Id      int    `json:"id"`
	GroupId int    `json:"group-id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Url     string `json:"url"`
}
