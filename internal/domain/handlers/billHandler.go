package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"psi-system.be.go.fiber/internal/domain/services"
	"psi-system.be.go.fiber/internal/domain/utils"
	"strings"
)

type BillToReceiveHandler struct {
	Service services.TransactionService
	Logger  *logrus.Logger
}

func NewBillToReceiveHandler(service services.TransactionService) *BillToReceiveHandler {
	return &BillToReceiveHandler{
		Service: service,
	}
}

func (h *BillToReceiveHandler) ListBillByType(c *fiber.Ctx) error {
	typeParam := c.Query("type")

	if typeParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "O parâmetro 'type' é obrigatório"})
	}

	var transactionType enums.TransactionType
	transactionType, err := utils.CastToTransactioTypeEnum(typeParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Não foi possível identificar o tipo de transação, por favor atualize a página e tente novamente"})
	}

	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Você não está logado no sistema, por favor faça o login novamente"})
	}

	bills, err := h.Service.ListBillByType(psychologistID, transactionType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao listar faturas a receber"})
	}

	return c.Status(http.StatusOK).JSON(bills)
}

func (h *BillToReceiveHandler) CreateBill(c *fiber.Ctx) error {
	var billDTO cashflow.BillDTO
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Você não está logado no sistema, por favor faça o login novamente"})
	}
	if err := c.BodyParser(&billDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ocorreu uma falha inesperada, por favor tente novamente mais tarde"})
	}

	if err := h.Service.CreateBill(&billDTO, psychologistID); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateBill",
		}).Error("Error saving bill to receive: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving bill to receive"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Bill to receive created successfully!"})
}

func (h *BillToReceiveHandler) UpdateBill(c *fiber.Ctx) error {
	var billDTO cashflow.BillDTO

	if err := c.BodyParser(&billDTO); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "UpdateBill",
		}).Error("Failed to read body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input format"})
	}

	if billDTO.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID must be provided"})
	}

	if err := h.Service.UpdateBill(&billDTO); err != nil {

		if strings.Contains(err.Error(), "Bill not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Não foi possível encontrar a fatura, verifique se as informações fornecidas estão corretas e tente novamente"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor tente novamente mais tarde"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Atualizado com sucesso"})
}

func (h *BillToReceiveHandler) ConfirmPaymentBill(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Fatura não encontrada. Verifique se a mesma encontra-se no portal e tente novamente"})
	}
	bill, err := h.Service.GetByID(uint64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor tente novamente mais tarde"})
	}

	bill.Status = enums.PAID.String()

	if err := h.Service.StatusUpdateBill(bill); err != nil {

		if strings.Contains(err.Error(), "Bill not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Não foi possível encontrar a fatura, verifique se as informações fornecidas estão corretas e tente novamente"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor tente novamente mais tarde"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Atualizado com sucesso"})
}

func (h *BillToReceiveHandler) RemoveConfirmationPaymentBill(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Fatura não encontrada. Verifique se a mesma encontra-se no portal e tente novamente"})
	}
	bill, err := h.Service.GetByID(uint64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor tente novamente mais tarde"})
	}

	bill.Status = enums.PAID.String()

	if err := h.Service.StatusUpdateBill(bill); err != nil {

		if strings.Contains(err.Error(), "Bill not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Não foi possível encontrar a fatura, verifique se as informações fornecidas estão corretas e tente novamente"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor tente novamente mais tarde"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Atualizado com sucesso"})
}

func (h *BillToReceiveHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "GetByID",
		}).Error("Error parsing ID: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	bill, err := h.Service.GetByID(uint64(id))
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "GetByID",
		}).Error("Error fetching bill to receive: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching bill to receive"})
	}

	if bill == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bill to receive not found"})
	}

	return c.Status(http.StatusOK).JSON(bill)
}

func (h *BillToReceiveHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "Delete",
		}).Error("Error parsing ID: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Registro não encontrado. Atualize a página e verifique se o registro ainda encontra-se disponível."})
	}

	if err := h.Service.Delete(uint64(id)); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "Delete",
		}).Error("Error deleting bill to receive: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ocorreu um erro na tentativa de deletar a fatura, por favor tente novamente mais tarde."})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Bill to receive deleted successfully"})
}
