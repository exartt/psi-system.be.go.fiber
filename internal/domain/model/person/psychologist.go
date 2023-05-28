package person

import "psi-system.be.go.fiber/internal/domain/model/appointment"

type Psychologist struct {
	ID           uint `gorm:"primary_key"`
	PersonID     uint `gorm:"not null"`
	Access       int
	TenantID     uint                      `gorm:"not null"`
	Password     string                    `gorm:"not null"`
	Appointments []appointment.Appointment `gorm:"ForeignKey:PsychologistID"`
}
