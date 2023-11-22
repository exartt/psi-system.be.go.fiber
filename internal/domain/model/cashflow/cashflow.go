package cashflow

import (
	"time"
)

type CashFlow struct {
	ID              uint      `gorm:"primary_key;column:id_fluxo_caixa"`
	PsychologistID  uint      `gorm:"column:id_psicologo"`
	PatientID       uint      `gorm:"column:id_paciente;index"`
	AppointmentID   uint      `gorm:"column:id_appointment;index"`
	TransactionType string    `gorm:"column:tipo_transacao"`
	Value           float64   `gorm:"type:float(8);column:flu_valor"`
	Description     string    `gorm:"column:descricao"`
	RecordDate      time.Time `gorm:"type:date;column:flu_data_registro"`
	Status          string    `gorm:"column:status"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

type BillToReceiveMonthly struct {
	PatientID      uint      `gorm:"column:id_paciente;"`
	PsychologistID uint      `gorm:"column:id_psicologo"`
	Value          float64   `gorm:"column:flu_valor"`
	Description    string    `gorm:"column:descricao"`
	TipoTransacao  string    `gorm:"column:tipo_transacao"`
	Status         string    `gorm:"column:status"`
	RecordDate     time.Time `gorm:"column:flu_data_registro"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

type DTO struct {
	ID              uint
	PatientID       uint
	TransactionType string
	Value           float64
	Description     string
	RecordDate      time.Time
	PatientName     string
	UpdatedAt       time.Time
}

type Table struct {
	Description     string    `gorm:"column:descricao"`
	ID              uint      `gorm:"column:id_fluxo_caixa"`
	TransactionType string    `gorm:"column:tipo_transacao"`
	Value           float64   `gorm:"column:flu_valor"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}
