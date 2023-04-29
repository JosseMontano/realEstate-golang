package controllers

import (
	"strconv"

	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/gofiber/fiber/v2"
)

func AddFavorite(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	userId, _ := strconv.ParseUint(data["user_id"], 10, 32)
	REId, _ := strconv.ParseUint(data["real_estate_id"], 10, 32)

	favoriteRE := models.FavoriteRealEstate{
		RealEstateId: uint(REId),
		UserId:       uint(userId),
	}

	database.DB.Create(&favoriteRE)
	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Guardado correctamente",
		"data":    favoriteRE,
	})
}
