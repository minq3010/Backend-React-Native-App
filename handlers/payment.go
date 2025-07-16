package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minq3010/Backend-React-Native-App/models"
	"github.com/minq3010/Backend-React-Native-App/repositories"
	"github.com/minq3010/Backend-React-Native-App/utils"
)

type PaymentHandler struct {
	Repo *repositories.PaymentRepository
}

func NewPaymentHandler(repo *repositories.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{Repo: repo}
}

func (h *PaymentHandler) CreateVnpayCheckout(c *fiber.Ctx) error {
	type Req struct {
		TicketID string `json:"ticket_id"`
		Amount   int    `json:"amount"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	orderID := req.TicketID + time.Now().Format("20060102150405")
	url, err := utils.CreateVnpayURL(orderID, req.Amount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	payment := &models.Payment{
		OrderID:       orderID,
		TicketID:      req.TicketID,
		Amount:        req.Amount,
		Status:        "pending",
		PaymentMethod: "vnpay",
	}
	if err := h.Repo.Create(context.Background(), payment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "create failed"})
	}

	return c.JSON(fiber.Map{"url": url})
}
