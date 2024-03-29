package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func FiberMiddleware(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowOrigins: "https://outclass.azlir.my.id, https://outclass-dev.netlify.app, https://master--outclass-dev.netlify.app, https://localhost:5173, http://localhost:5173",
		}),
		logger.New(),
	)
}
