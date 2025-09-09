package routes

import (
	"task-manager-app/controllers"
	"task-manager-app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupProjectRoutes(app *fiber.App){
	projectGroup := app.Group("/projects")

	projectGroup.Use(middlewares.AuthMiddleWare)
	
	projectGroup.Post("/", controllers.CreateProject)
	projectGroup.Get("/", controllers.GetProjects)
	projectGroup.Get("/:id", controllers.GetProject)
	projectGroup.Put("/:id", controllers.UpdateProject)
	projectGroup.Delete("/:id", controllers.DeleteProject)
}