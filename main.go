package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/s6352410016/go-fiber-gorm-api-auth-jwt-mysql/config"
	"github.com/s6352410016/go-fiber-gorm-api-auth-jwt-mysql/database"
	"github.com/s6352410016/go-fiber-gorm-api-auth-jwt-mysql/routes"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	app := fiber.New()
	routes.SetUpRoutes(app)

	app.Listen(":8080")
}
