package main

import (
	"fmt"
	"log"
	"os"
	"outclass-api/app/configs"
	"outclass-api/app/middleware"
	"outclass-api/app/routes"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := configs.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	routes.PrivateRoutes(app)
	routes.PublicRoutes(app)

	url := fmt.Sprintf(
		"%s:%s",
		os.Getenv("SERVER_HOST"),
		os.Getenv("SERVER_PORT"),
	)
	log.Fatal(app.Listen(url))
}