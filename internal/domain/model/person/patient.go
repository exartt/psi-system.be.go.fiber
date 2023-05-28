package person

import "psi-system.be.go.fiber/internal/domain/model/appointment"

type Patient struct {
	ID             uint `gorm:"primary_key"`
	PersonID       uint `gorm:"not null"`
	PsychologistID uint `gorm:"not null"`
	IsPlan         bool
	SessionPrice   float64
	ConversionType string
	Appointments   []appointment.Appointment `gorm:"foreignKey:PatientID"`
}
