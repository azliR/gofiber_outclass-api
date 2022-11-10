package routes

import (
	"outclass-api/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Post("/user/sign/in", handlers.UserSignIn)
	route.Post("/user/sign/up", handlers.UserSignUp)
}
