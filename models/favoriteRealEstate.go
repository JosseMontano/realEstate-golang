package models

type FavoriteRealEstate struct {
	Id           uint `json:"id"`
	RealEstateId uint `json:"real_estate_id"`
	UserId       uint `json:"user_id"`
}
