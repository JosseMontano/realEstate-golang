package controllers

import (
	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user []models.User
	database.DB.Where("id = ?", id).Find(&user)
	c.Status(200)
	return c.JSON(user)
}
