package models

type Photo struct {
	Id       uint   `json:"id"`
	Url      string `json:"url"`
	PublicId string `json:"public_id"`
}
