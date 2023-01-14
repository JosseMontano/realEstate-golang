package routes

import (
	"github.com/JosseMontano/estateInTheCloud/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/api/signup", controllers.Register)
	app.Post("/api/signin", controllers.SingIn)
	app.Get("/api/me", controllers.User)
	app.Post("/api/logout", controllers.Logout)
}
