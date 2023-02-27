package controllers

import (
	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/gofiber/fiber/v2"
)

func AllTypeRE(c *fiber.Ctx) error {
	var typeRealEstate []models.TypeRealEstate
	database.DB.Find(&typeRealEstate)
	c.Status(200)
	return c.JSON(typeRealEstate)
}
