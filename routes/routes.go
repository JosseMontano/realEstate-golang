package routes

import (
	"github.com/JosseMontano/estateInTheCloud/controllers"
	"github.com/JosseMontano/estateInTheCloud/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/api/signup", controllers.Register)
	app.Post("/api/signin", controllers.SingIn)
	app.Post("/api/recuperateAccount", controllers.SendCodeToGmail)

	app.Use(middleware.ValidateJwt)

	//Sign
	app.Get("/api/me", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	//========== REAL ESTATE ==========
	app.Get("/api/realEstate", controllers.AllRE)
	app.Get("/api/estateMostRecent", controllers.MostRecentRE)
	app.Get("/api/estateRecommendedByUser", controllers.UserRecommend)
	app.Post("/api/realEstate", controllers.CreateRE)

}
