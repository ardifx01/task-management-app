package main

import (
	"log"
	"task-manager-app/config"
	"task-manager-app/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectDB()

	app := fiber.New()
	
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	routes.SetupAuthRoutes(app)
	routes.SetupProjectRoutes(app)
	routes.SetupTaskRoutes(app)

	log.Fatal(app.Listen(":3000"))
}