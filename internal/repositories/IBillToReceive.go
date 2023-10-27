package repositories

import (
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
)

type BillToReceiveRepository interface {
	CreateBillToReceive(billToReceive *cashflow.CashFlow) error
	ListBillToReceive() ([]*cashflow.CashFlow, error)
	UpdateBillToReceive(billToReceive *cashflow.CashFlow) error
	GetByID(id uint64) (*cashflow.CashFlow, error)
}

func (r *cashFlowRepository) CreateBillToReceive(billToReceive *cashflow.CashFlow) error {
	result := r.db.Create(billToReceive)
	return result.Error
}

func (r *cashFlowRepository) ListBillToReceive() ([]*cashflow.CashFlow, error) {
	var bills []*cashflow.CashFlow
	result := r.db.Where("transaction_type = ?", enums.RECEIVABLE.String()).Find(&bills)
	return bills, result.Error
}

func (r *cashFlowRepository) UpdateBillToReceive(billToReceive *cashflow.CashFlow) error {
	result := r.db.Save(billToReceive)
	return result.Error
}
