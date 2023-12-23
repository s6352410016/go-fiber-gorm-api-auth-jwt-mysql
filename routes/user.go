package routes

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/s6352410016/go-fiber-gorm-api-auth-jwt-mysql/handlers"
)

func SetUpRoutes(app *fiber.App) {
	user := app.Group("/api")

	user.Use("/profile", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("AT_SECRET"))},
	}))
	user.Use("/refresh", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("RT_SECRET"))},
	}))

	user.Post("/signup", handlers.SignUp)
	user.Post("/signin", handlers.SignIn)
	user.Get("/profile", handlers.Profile)
	user.Post("/refresh", handlers.Refresh)
}
