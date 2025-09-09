package routes

import (
	"task-manager-app/controllers"
	"task-manager-app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(app *fiber.App){
	taskGroup := app.Group("/")
	taskGroup.Use(middlewares.AuthMiddleWare)

	taskGroup.Get("/projects/:id/tasks", controllers.GetTasks)
	taskGroup.Post("/projects/:id/tasks", controllers.CreateTask)

	taskGroup.Get("/tasks/:id", controllers.GetTask)
	taskGroup.Put("/tasks/:id", controllers.UpdateTask)
	taskGroup.Delete("/tasks/:id", controllers.DeleteTask)
	taskGroup.Get("/tasks", controllers.FilterTask)
}