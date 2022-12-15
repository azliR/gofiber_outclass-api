package middleware

import (
	"os"
	"outclass-api/app/commons"

	"github.com/gofiber/fiber/v2"

	jwtMiddleware "github.com/gofiber/jwt/v3"
)

func JWTProtected() func(*fiber.Ctx) error {
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey:   "jwt",
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
		Success: false,
		Message: err.Error(),
	})
}
