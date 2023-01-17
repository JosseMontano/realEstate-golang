package controllers

import (
	"fmt"
	"strconv"

	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateStructRE(realEstate models.RealEstate) []*models.ErrorResponseRE {
	var errors []*models.ErrorResponseRE
	err := validate.Struct(realEstate)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponseRE
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func AllRE(c *fiber.Ctx) error {
	var realEstate []models.RealEstate
	database.DB.Preload("User").Preload("TypeRealState").Preload("Photos").Find(&realEstate)
	return c.JSON(realEstate)
}

func MostRecentRE(c *fiber.Ctx) error {
	var realEstate []models.RealEstate
	database.DB.Limit(10).Preload("User").
		Preload("TypeRealState").Preload("Photos").Order("id desc").Find(&realEstate)
	return c.JSON(realEstate)
}

type RealState struct {
}

func UserRecommend(c *fiber.Ctx) error {
	var user []models.User

	//this working
	/* 	database.DB.Debug().Table("(?) as u",
	database.DB.Debug().Joins("JOIN users on real_estates.user_id=users.id").
		Order("users.qualification desc").Preload("User").Find(&realEstate)).
	Distinct("user_id").Preload("User").Find(&realEstate) */

	database.DB.Order("qualification desc").Find(&user)

	return c.JSON(user)
}

func CreateRE(c *fiber.Ctx) error {
	/* 	var realEstate models.RealEstate */
	var realEstateDto fiber.Map
	if err := c.BodyParser(&realEstateDto); err != nil {
		return err
	}

	/* errors := ValidateStructRE(realEstateDto)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors,
		})
	} */

	list := realEstateDto["photos"].([]interface{})
	photos := make([]models.Photo, len(list))

	for i, photosId := range list {
		id, _ := strconv.Atoi(photosId.(string))
		photos[i] = models.Photo{
			Id: uint(id),
		}
	}

	fmt.Println(realEstateDto)

	realEstate := models.RealEstate{
		Title:            realEstateDto["title"].(string),
		Description:      realEstateDto["description"].(string),
		AmountBedroom:    int(realEstateDto["amount_bedroom"].(float64)),
		Price:            int(realEstateDto["price"].(float64)),
		AmountBathroom:   int(realEstateDto["amount_bathroom"].(float64)),
		SquareMeter:      int(realEstateDto["square_meter"].(float64)),
		UserId:           int(realEstateDto["user_id"].(float64)),
		TypeRealEstateId: int(realEstateDto["type_real_estate_id"].(float64)),
		Photos:           photos,
	}

	database.DB.Create(&realEstate)
	database.DB.Model(&realEstate).Association("User").Find(&realEstate.User)
	database.DB.Model(&realEstate).Association("TypeRealState").Find(&realEstate.TypeRealState)
	database.DB.Model(&realEstate).Association("Photos").Find(&realEstate.Photos)
	return c.JSON(realEstate)
}
