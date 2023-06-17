package cashflow

import "time"

type BillToPay struct {
	CashFlowID uint      `gorm:"primary_key;column:id_fluxo_caixa"`
	DateToPay  time.Time `gorm:"type:date;column:data_a_pagar"`
	Status     bool      `gorm:"column:status"`
}
