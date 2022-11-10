package routes

import (
	"outclass-api/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Post("/user/sign/in", controllers.UserSignIn)
	route.Post("/user/sign/up", controllers.UserSignUp)
}
