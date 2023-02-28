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

type AllREResult struct {
	IdRealEstate      int    `json:"id_real_estate"`
	IdRealEstatePhoto int    `json:"id_real_estate_photo"`
	IdPhoto           int    `json:"id_photo"`
	Url               string `json:"url"`
	PublicId          string `json:"public_id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Email             string `json:"email"`
	IdUser            int    `json:"id_user"`
}

const from = "from real_estates_photos rp , photos p, real_estates re, users u"

const where = "where rp.photo_id = p.id and rp.real_estate_id = re.id and re.user_id = u.id"

const selectQuery = `select DISTINCT on (re.id) re.id as id_real_estate, rp.id as id_real_estate_photo,
p.id as id_photo, p.url, p.public_id, re.title,
re.description, u.email, u.id as id_user`

const query = selectQuery + " " + from + " " + where + " " + "and re.available=true ORDER BY re.id"

func AllRE(c *fiber.Ctx) error {
	var realEstate []AllREResult
	database.DB.Debug().Raw(query).Scan(&realEstate)
	return c.JSON(realEstate)
}

func MostRecentRE(c *fiber.Ctx) error {
	var realEstate []AllREResult
	database.DB.Debug().Raw(query + " " + "desc limit 8").Scan(&realEstate)
	return c.JSON(realEstate)
}

func UserRecommend(c *fiber.Ctx) error {
	var realEstate []AllREResult
	database.DB.Raw(`SELECT * 
	FROM(SELECT DISTINCT on (u.email) re.id as id_real_estate, rp.id as id_real_estate_photo,
	p.id as id_photo,  p.url, 
	p.public_id, re.title, re.description, u.email, u.id as id_user, u.qualification` +
		" " + from + " " + where + " " +
		`and re.available=true ORDER BY u.email DESC) users ORDER BY users.qualification desc`).Scan(&realEstate)
	return c.JSON(realEstate)
}

func RealEstate(c *fiber.Ctx) error {
	id := c.Params("id")
	var realEstate []models.RealEstate
	database.DB.Where("id=?", id).Preload("User").Preload("Photos").Find(&realEstate)
	c.Status(200)
	return c.JSON(realEstate)
}

func RealEstateByType(c *fiber.Ctx) error {
	typeRE := c.Params("type")
	var realEstate []AllREResult
	database.DB.Raw(selectQuery + " " + from + ", " + "type_real_estates tre" + " " + where + " " +
		`and re.available=true and re.type_real_estate_id = tre.id and tre.name='` + typeRE + `' ORDER BY re.id`).Scan(&realEstate)
	c.Status(200)
	return c.JSON(realEstate)
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
	database.DB.Model(&realEstate).Association("TypeRealState").Find(&realEstate.TypeRealEstate)
	database.DB.Model(&realEstate).Association("Photos").Find(&realEstate.Photos)
	return c.JSON(realEstate)
}
