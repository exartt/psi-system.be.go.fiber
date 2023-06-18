package cashflow

import (
	"psi-system.be.go.fiber/internal/domain/model/appointment"
	"time"
)

type CashFlow struct {
	ID              uint                    `gorm:"primary_key;column:id_fluxo_caixa"`
	PsychologistID  uint                    `gorm:"column:id_psicologo"`
	PatientId       uint                    `gorm:"column:id_paciente;index"`
	AppointmentID   uint                    `gorm:"column:id_appointment;index"`
	Appointment     appointment.Appointment `gorm:"foreignKey:AppointmentID"`
	TransactionType string                  `gorm:"column:tipo_transacao"`
	Value           float64                 `gorm:"type:float(8);column:flu_valor"`
	Description     string                  `gorm:"column:descricao"`
	RecordDate      time.Time               `gorm:"type:date;column:flu_data_registro"`
	Status          string                  `gorm:"column:status"`
	CreatedAt       time.Time               `gorm:"column:created_at"`
	UpdatedAt       time.Time               `gorm:"column:updated_at"`
	CreatedBy       string                  `gorm:"column:created_by"`
}
