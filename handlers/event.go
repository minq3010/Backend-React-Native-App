package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minq3010/Backend-React-Native-App/models"
)

type EventHandler struct {
	repository models.EventRepository
}

func 

func NewEventHandler( router fiber.Router, repository models.EventRepository) {
	handler := &EventHandler {
		repository: repository,
	}
router.Get("/", handler.GetMany)
router.Post("/", handler.CreateOne)
router.Get("/:eventId", handler.GetOne)
}
