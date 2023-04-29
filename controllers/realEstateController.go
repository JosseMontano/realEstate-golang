package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/helper"
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
	//* ===== SAVE IMG IN CLOUDINARY =====
	formHeader, err := c.FormFile("image")

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"StatusCode": http.StatusInternalServerError,
				"Message":    "error",
				"Data":       "Select a file to upload",
			})
	}

	formFile, err := formHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"StatusCode": http.StatusInternalServerError,
				"Message":    "error",
				"Data":       err.Error(),
			})
	}

	url, state := helper.Upload(formFile)

	if state {
		return c.JSON(fiber.Map{
			"message": url,
		})
	}

	return c.JSON(url)

}

type TypeGetTypeRE struct {
	TypeRealEstateId int `json:"type_real_estate_id"`
}

type TypeGetMinMaxRE struct {
	MaxSquareMeter   int `json:"max_square_meter"`
	MinSquareMeter   int `json:"min_square_meter"`
	MaxAmountBedroom int `json:"max_amount_bedroom"`
	MinAmountBedroom int `json:"min_amount_bedroom"`
}

func FilterIntelligente(c *fiber.Ctx) error {
	userId := c.Params("user_id")

	//===== get the types of real estate =====
	var typesRealEstate []TypeGetTypeRE
	database.DB.Raw(`select DISTINCT on (RE.type_real_estate_id) RE.type_real_estate_id
	from favorite_real_estates fav, real_estates RE
	where fav.real_estate_id = RE.id and Fav.user_id=` + userId).Scan(&typesRealEstate)

	//===== get the min and max =====
	var MinMaxRE TypeGetMinMaxRE
	database.DB.Raw(`select max(RE.square_meter) as max_square_meter, min(RE.square_meter) as min_square_meter,
	max(RE.amount_bedroom) as max_amount_bedroom, min(RE.amount_bedroom) as min_amount_bedroom
	from favorite_real_estates fav, real_estates RE
	where fav.real_estate_id = RE.id and Fav.user_id=` + userId).Scan(&MinMaxRE)

	//===== FILTER IA =====
	var realEstate []AllREResult

	var conditions []string
	for _, t := range typesRealEstate {
		conditions = append(conditions, fmt.Sprintf("type_real_estate_id = %d", t.TypeRealEstateId))
	}

	queryWhere := " and square_meter >=" + strconv.Itoa(MinMaxRE.MinSquareMeter) + " and square_meter <=" +
		strconv.Itoa(MinMaxRE.MaxSquareMeter) +
		" and amount_bedroom <=" + strconv.Itoa(MinMaxRE.MinAmountBedroom) +
		" and amount_bedroom >=" + strconv.Itoa(MinMaxRE.MaxAmountBedroom)

	query := fmt.Sprintf(`
	SELECT  re.id as id_real_estate, re.title, re.description,
        rp.id as id_real_estate_photo,
        p.id as id_photo, p.url, p.public_id, re.title,
        re.description, u.email, u.id as id_user
FROM real_estates re, real_estates_photos rp , photos p, users u
WHERE rp.photo_id = p.id and rp.real_estate_id = re.id and re.user_id = u.id 
  and re.id IN (SELECT DISTINCT ON (re.id) re.id as id_real_estate
                FROM real_estates re, users u
                WHERE re.user_id = u.id and %s
	`, strings.Join(conditions, " OR "))
	database.DB.Debug().Raw(query + queryWhere + ")").Scan(&realEstate)

	return c.JSON(realEstate)
}

/*
	query := fmt.Sprintf(`SELECT DISTINCT on (re.id) re.id as id_real_estate, re.title, re.description,
	rp.id as id_real_estate_photo,
	p.id as id_photo, p.url, p.public_id, re.title,
	re.description, u.email, u.id as id_user
	FROM real_estates re, real_estates_photos rp , photos p, users u
	WHERE rp.photo_id = p.id and rp.real_estate_id = re.id and re.user_id = u.id and %s
	`, strings.Join(conditions, " OR "))
	database.DB.Debug().Raw(query + queryWhere).Scan(&realEstate)

*/
