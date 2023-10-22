package googleCalendar

import "time"

type Token struct {
	ID                 uint      `gorm:"primary_key"`
	PsyID              uint      `gorm:"not null;foreignKey:ID;references:Psychologists"`
	Token              string    `gorm:"type:text;not null"`
	ExpirationDateTime time.Time `gorm:"not null"`
	TenantID           uint      `gorm:"not null"`
}
