package googleCalendar

import "time"

type TokenDTO struct {
	ID                 uint      `json:"id"`
	PsyId              uint      `json:"psychologistId"`
	Token              string    `json:"token"`
	ExpirationDateTime time.Time `json:"expirationDateTime"`
	TenantID           uint      `json:"tenantId"`
}
