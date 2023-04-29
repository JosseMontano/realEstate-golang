package main

import (
	"log"
	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/routes"
	"github.com/JosseMontano/estateInTheCloud/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.Connect()

	app := fiber.New()
	//https://realestate-c70dc.web.app
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))

	routes.Setup(app)

	port := utils.DotEnvVariable("PORT")

	/* log.Fatal(app.Listen("0.0.0.0:" + port)) */
	log.Fatal(app.Listen(port))
}
