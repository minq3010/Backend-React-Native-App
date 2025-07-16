package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/minq3010/Backend-React-Native-App/models"
	"github.com/minq3010/Backend-React-Native-App/utils"
)

type PaymentHandler struct {
	PaymentRepo models.PaymentRepository
	EventRepo   models.EventRepository
	TicketRepo  models.TicketRepository
}

// POST /payment/vnpay
func (h *PaymentHandler) CreateVnpayCheckout(c *fiber.Ctx) error {
	type Body struct {
		EventID uint `json:"eventId"`
	}
	var body Body

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	userId := uint(c.Locals("userId").(float64))
	ctx := context.Background()

	event, err := h.EventRepo.GetOne(ctx, body.EventID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "event not found",
		})
	}

	orderID := utils.GenerateOrderID(userId, event.ID)

	payment := &models.Payment{
		OrderID:       orderID,
		TicketID:      "", // chưa có ticket
		EventID:       event.ID,
		UserID:        userId,
		Amount:        event.Price, // ép kiểu nếu cần
		Status:        "pending",
		PaymentMethod: "vnpay",
	}

	if err := h.PaymentRepo.Create(ctx, payment); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create payment",
		})
	}

	payURL, err := utils.CreateVnpayURL(orderID, payment.Amount)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate payment URL",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"url":    payURL,
	})
}

// GET /payment/vnpay-return
func (h *PaymentHandler) HandleVnpayReturn(c *fiber.Ctx) error {
	query := c.Queries()
	orderID := query["vnp_TxnRef"]
	responseCode := query["vnp_ResponseCode"]
	//signature := query["vnp_SecureHash"]

	// 👉 In log để kiểm tra callback và chữ ký
	fmt.Println("🔍 vnp_ResponseCode:", responseCode)
	fmt.Println("🔍 Raw query string:", c.Context().QueryArgs().String())
	fmt.Println("🔍 Parsed Queries():", query)

	// 1. Xác thực chữ ký
	// if !utils.VerifyVnpaySignature(query, signature) {
	// 	return c.Status(http.StatusBadRequest).SendString("❌ Sai chữ ký, không hợp lệ")
	// }

	// 2. Kiểm tra mã phản hồi từ VNPAY
	if responseCode != "00" {
		return c.SendString("❌ Thanh toán bị từ chối hoặc thất bại")
	}

	ctx := context.Background()

	// 3. Tìm đơn thanh toán trong DB
	payment, err := h.PaymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("❌ Không tìm thấy đơn thanh toán")
	}

	if payment.Status == "success" {
		return c.SendString("✅ Đơn đã xử lý trước đó")
	}

	// 4. Cập nhật trạng thái đơn thanh toán
	if err := h.PaymentRepo.UpdateStatus(ctx, orderID, "success"); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("❌ Cập nhật trạng thái lỗi")
	}

	// 5. Tạo ticket cho user
	ticket := &models.Ticket{
		UserID:  payment.UserID,
		EventID: payment.EventID,
		Price:   payment.Amount,
	}

	createdTicket, err := h.TicketRepo.CreateOne(ctx, payment.UserID, ticket)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("❌ Tạo vé thất bại")
	}

	// 6. ✅ Cập nhật TicketID vào đơn thanh toán
	err = h.PaymentRepo.UpdateTicketID(ctx, orderID, fmt.Sprintf("%d", createdTicket.ID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("❌ Gán TicketID vào Payment thất bại")
	}

	return c.SendString("✅ Thanh toán thành công, vé đã được tạo!")
}

func NewPaymentHandler(router fiber.Router, pRepo models.PaymentRepository, eRepo models.EventRepository, tRepo models.TicketRepository) {
	handler := &PaymentHandler{
		PaymentRepo: pRepo,
		EventRepo:   eRepo,
		TicketRepo:  tRepo,
	}

	router.Post("/vnpay", handler.CreateVnpayCheckout)
	router.Get("/vnpay-return", handler.HandleVnpayReturn)
}
