package models

//idUser idTypeRealEstate
type RealEstate struct {
	Id               uint          `json:"id"`
	Title            string        `json:"title" validate:"required"`
	Description      string        `json:"description" validate:"required"`
	Available        bool          `json:"available"`
	AmountBedroom    int           `json:"amount_bedroom" validate:"required,number"`
	Price            int           `json:"price" validate:"required,number"`
	AmountBathroom   int           `json:"amount_bathroom" validate:"required,number"`
	SquareMeter      int           `json:"square_meter" validate:"required,number"`
	UserId           int           `json:"user_id"`
	User             User          `json:"user" gorm:"foreignKey:UserId"`
	TypeRealEstateId int           `json:"type_real_estate_id"`
	TypeRealState    TypeRealState `json:"type_real_state" gorm:"foreignKey:TypeRealEstateId"`
	Photos           []Photo       `json:"photos" gorm:"many2many:real_estates_photos"`
}

type ErrorResponseRE struct {
	FailedField string
	Tag         string
	Value       string
}
