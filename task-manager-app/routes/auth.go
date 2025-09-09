package routes

import (
	"task-manager-app/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App){ 
	authGroup := app.Group("/auth")
	authGroup.Post("/register", controllers.Register)
	authGroup.Post("/login", controllers.Login)
}