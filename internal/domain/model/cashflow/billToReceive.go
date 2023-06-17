package cashflow

import "time"

type BillToReceive struct {
	CashFlowID    uint      `gorm:"primary_key;column:id_fluxo_caixa"`
	DateToReceive time.Time `gorm:"type:date;column:data_a_receber"`
	Status        bool      `gorm:"column:status"`
}
