package models

type Question struct {
	Id       uint   `json:"id"`
	Question string `json:"question" validate:"required,min=3,max=32"`
}
