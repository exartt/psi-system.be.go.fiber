package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"psi-system.be.go.fiber/internal/domain/services"
	"strings"
)

type BillToReceiveHandler struct {
	Service services.BillToReceiveService
	Logger  *logrus.Logger
}

func NewBillToReceiveHandler(service services.BillToReceiveService) *BillToReceiveHandler {
	return &BillToReceiveHandler{
		Service: service,
	}
}

func (h *BillToReceiveHandler) ListBillToReceive(c *fiber.Ctx) error {
	bills, err := h.Service.ListBillToReceive()
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "ListBillToReceive",
		}).Error("Error listing bills to receive: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error listing bills to receive"})
	}
	return c.Status(http.StatusOK).JSON(bills)
}

func (h *BillToReceiveHandler) CreateBillToReceive(c *fiber.Ctx) error {
	var billDTO cashflow.BillToPayDTO

	if err := c.BodyParser(&billDTO); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateBillToReceive",
		}).Error("Failed to read body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
	}

	if err := h.Service.CreateBillToReceive(&billDTO); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateBillToReceive",
		}).Error("Error saving bill to receive: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving bill to receive"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Bill to receive created successfully!"})
}

func (h *BillToReceiveHandler) UpdateBillToReceive(c *fiber.Ctx) error {
	var billDTO cashflow.BillToPayDTO

	if err := c.BodyParser(&billDTO); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "UpdateBillToReceive",
		}).Error("Failed to read body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input format"})
	}

	if billDTO.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID must be provided"})
	}

	if err := h.Service.UpdateBillToReceive(&billDTO); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "UpdateBillToReceive",
		}).Error("Error updating bill to receive: ", err)

		if strings.Contains(err.Error(), "Bill not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bill not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating bill to receive"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Bill to receive updated successfully"})
}
