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

	serverScheme := os.Getenv("SERVER_SCHEME")
	serverPort := os.Getenv("SERVER_PORT")
	if serverScheme == "http" {
		log.Fatal(app.Listen(":" + serverPort))
	} else if serverScheme == "https" {
		log.Fatal(app.ListenTLS(":" + serverPort, "./cert.pem", "./privkey.pem"))
	}
}
