package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func FiberMiddleware(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowOrigins: "https://outclass-dev.netlify.app, http://localhost:5000",
		}),
		logger.New(),
	)
}
