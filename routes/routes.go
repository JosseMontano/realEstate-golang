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

	//======================= AUTH =======================

	app.Post("/api/signup", controllers.Register)
	app.Post("/api/signin", controllers.SingIn)
	app.Post("/api/sendCodeGmail", controllers.SendCodeToGmail)
	app.Post("/api/changePassword", controllers.ChangePassword)

	app.Use(middleware.ValidateJwt)

	//======================= USER =======================
	app.Get("/api/getUserComplete/:id", controllers.GetUser)

	//======================= AUTH =======================
	app.Get("/api/me", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	//======================= REAL ESTATE =======================
	app.Get("/api/allRealEstate", controllers.AllRE)
	app.Get("/api/typeRealEstate", controllers.AllTypeRE)
	app.Get("/api/estateMostRecent", controllers.MostRecentRE)
	app.Get("/api/estateRecommendedByUser", controllers.UserRecommend)
	app.Get("/api/realEstate/:id", controllers.RealEstate)



	app.Post("/api/realEstate", controllers.CreateRE)

}
