package main

import (
	"fmt"

	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/routes"
	"github.com/JosseMontano/estateInTheCloud/utils"
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

	port := utils.DotEnvVariable("PORT_SERVER")

	fmt.Print(port)
	app.Listen(port)
}
