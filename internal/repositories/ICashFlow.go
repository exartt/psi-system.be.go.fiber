package repositories

import (
	"errors"
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
)

type CashFlowRepository interface {
	CreateBillToReceive(billToReceive *cashflow.CashFlow) error
	GetByID(id uint64) (*cashflow.CashFlow, error)
	ListBillToReceive() ([]*cashflow.CashFlow, error)
	UpdateBillToReceive(billToReceive *cashflow.CashFlow) error
}

type cashFlowRepository struct {
	db *gorm.DB
}

func NewCashFlowRepository(db *gorm.DB) CashFlowRepository {
	return &cashFlowRepository{db: db}
}

func (r *cashFlowRepository) GetByID(id uint64) (*cashflow.CashFlow, error) {
	var bill cashflow.CashFlow
	result := r.db.First(&bill, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &bill, nil
}
