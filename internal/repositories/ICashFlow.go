package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
)

type CashFlowRepository interface {
	CreateBillToReceive(billToReceive *cashflow.CashFlow) error
}

type cashFlowRepository struct {
	db *gorm.DB
}

func NewCashFlowRepository(db *gorm.DB) CashFlowRepository {
	return &cashFlowRepository{db: db}
}
