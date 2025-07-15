package handlers

import "github.com/minq3010/Backend-React-Native-App/repositories"

type PaymentHandler struct {
	Repo *repositories.PaymentRepository
}

func NewPaymentHandler(repo *repositories.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{
		Repo: repo,
	}
}

func (h *PaymentHandler) CreateVnpayCheckout(c *fiber.Ctx) error {
	type Req struct {
		TicketID string `json:"ticket_id"`
		Amount int `json:"amount`

	}
}