package appointment

import (
	"psi-system.be.go.fiber/internal/domain/enums"
	"time"
)

type Appointment struct {
	ID             uint `gorm:"primary_key"`
	PsychologistID uint
	PatientID      uint
	TenantID       uint `gorm:"not null"`
	CalendarID     string
	Start          time.Time
	End            time.Time
	Summary        string
	Description    string
	Location       string
	Status         enums.StatusAgendamento
	Notify         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
