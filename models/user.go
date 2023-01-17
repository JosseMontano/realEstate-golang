package models

type User struct {
	Id              uint   `json:"id"`
	UserName        string `json:"user_name" validate:"required,min=3,max=32"`
	CellphoneNumber string `json:"cellphone_number" validate:"required"`
	Email           string `json:"email" validate:"required,email" gorm:"unique"`
	Password        []byte `json:"password" validate:"required"`
	Qualification   int    `json:"qualification"`
	UrlPhoto        string `json:"url_photo"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}
