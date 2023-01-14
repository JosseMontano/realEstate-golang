package models

type User struct {
	Id              uint   `json:"id"`
	UserName        string `json:"user_name"`
	CellphoneNumber string `json:"cellphone_number"`
	Email           string `json:"email" gorm:"unique"`
	Password        []byte `json:"password"`
	Qualification   int    `json:"qualification"`
	UrlPhoto        string `json:"url_photo"`
}
