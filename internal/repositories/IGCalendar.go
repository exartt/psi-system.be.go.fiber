package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/googleCalendar"
	"psi-system.be.go.fiber/internal/domain/utils"
	"time"
)

type GoogleCalendarTokenRepository interface {
	Store(token string, expirationDateTime time.Time, psyID uint) error
	FindByPsyID(psyID uint, expirationDateTime time.Time) (*googleCalendar.TokenDTO, error)
	Update(token string, psyID uint) error
	Delete(psyID uint) error
}

type googleCalendarTokenRepository struct {
	db *gorm.DB
}

func NewGCalendarRepository(db *gorm.DB) GoogleCalendarTokenRepository {
	return &googleCalendarTokenRepository{db: db}
}

func (r *googleCalendarTokenRepository) Store(token string, expirationDateTime time.Time, psyID uint) error {
	googleCalendarToken := googleCalendar.TokenDTO{
		PsyId:              psyID,
		Token:              token,
		ExpirationDateTime: expirationDateTime,
		TenantID:           0,
	}
	return r.db.Create(&googleCalendarToken).Error
}

func (r *googleCalendarTokenRepository) FindByPsyID(psyID uint, expirationDateTime time.Time) (*googleCalendar.TokenDTO, error) {
	var googleCalendarToken googleCalendar.TokenDTO
	err := r.db.Where("psy_id = ? AND expiration_date_time > ?", psyID, expirationDateTime).First(&googleCalendarToken).Error
	return &googleCalendarToken, err
}

func (r *googleCalendarTokenRepository) Update(token string, psyID uint) error {
	return r.db.Model(&googleCalendar.TokenDTO{}).Where("psy_id = ?", psyID).Updates(googleCalendar.TokenDTO{
		Token:              token,
		ExpirationDateTime: utils.GetExpirationDate(),
	}).Error
}

// Delete It won't be used, just impl to complete the crud.
func (r *googleCalendarTokenRepository) Delete(psyID uint) error {
	return r.db.Where("psy_id = ?", psyID).Delete(&googleCalendar.TokenDTO{}).Error
}
