package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/minq3010/Backend-React-Native-App/config"
	"github.com/minq3010/Backend-React-Native-App/db"
	"github.com/minq3010/Backend-React-Native-App/handlers"
	"github.com/minq3010/Backend-React-Native-App/repositories"
)

func main() {
	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)


	app := fiber.New(fiber.Config{
		AppName: "TicketBooking",
		ServerHeader: "Fiber",
	})
	// repositories
	eventRepository := repositories.NewEventRepository(db)
	// routing
	server := app.Group("/api")
	// handler
	handlers.NewEventHandler(server.Group("/event"), eventRepository)


	app.Listen(fmt.Sprint(":" + envConfig.ServerPort))
}