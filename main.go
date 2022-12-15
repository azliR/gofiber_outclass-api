package main

import (
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

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	url := host + ":" + port
	log.Fatal(app.Listen(url))
}
