package cashflow

import "time"

type CashFlow struct {
	IDCashFlow     uint            `gorm:"primary_key;column:id_fluxo_caixa"`
	PsychologistID uint            `gorm:"column:id_psicologo"`
	PatientID      uint            `gorm:"column:id_paciente"`
	Value          float64         `gorm:"type:float(8);column:flu_valor"`
	RecordDate     time.Time       `gorm:"type:date;column:flu_data_registro"`
	BillsToPay     []BillToPay     `gorm:"foreignKey:CashFlowID;references:IDCashFlow"`
	BillsToReceive []BillToReceive `gorm:"foreignKey:CashFlowID;references:IDCashFlow"`
}
