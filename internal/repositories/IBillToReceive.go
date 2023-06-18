package repositories

import (
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
)

type BillToReceiveRepository interface {
	CreateBillToReceive(billToReceive *cashflow.CashFlow) error
}

func (r *cashFlowRepository) CreateBillToReceive(billToReceive *cashflow.CashFlow) error {
	result := r.db.Create(billToReceive)
	return result.Error
}
