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

type DTO struct {
	ID              uint
	PatientID       uint
	TransactionType string
	Value           float64
	Description     string
	RecordDate      time.Time
	PatientName     string
}

//Appointment     appointment.Appointment `gorm:"foreignKey:AppointmentID"`
