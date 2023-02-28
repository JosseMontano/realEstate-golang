package controllers

import (
	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/gofiber/fiber/v2"
)

func GetQuestion(c *fiber.Ctx) error {
	id := c.Params("idReal_Estate")
	var questions []models.Question
	database.DB.Raw(`select * from questions q where q.id NOT IN (
	select aq.id_question as idQuestion 
	from answers_questions aq, answers a 
	where aq.id_answer=a.id and a.id_real_estate=` + id + `)`).Scan(&questions)
	c.Status(200)
	return c.JSON(questions)
}

func CreateQuestion(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	questions := models.Question{
		Question: data["question"],
	}

	database.DB.Create(&questions)
	c.Status(200)
	return c.JSON(fiber.Map{
		"action": true,
	})
}


