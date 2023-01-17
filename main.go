package main

import (
	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":4000")
}
