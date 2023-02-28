package controllers

import (
	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetPhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	idNumber, _ := strconv.Atoi(id)
	var photos []models.Photo
	database.DB.Find(&photos, []int{idNumber})
	return c.JSON(photos)
}