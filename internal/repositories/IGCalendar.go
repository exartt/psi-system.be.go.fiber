package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/googleCalendar"
	"psi-system.be.go.fiber/internal/domain/utils"
	"time"
)

type GoogleCalendarTokenRepository interface {
	Store(token string, expirationDateTime time.Time, psyID uint) error
	FindByPsyID(psyID uint) (*googleCalendar.TokenDTO, error)
	Update(token string, psyID uint) error
	Delete(psyID uint) error
	FindTokenByPsychologistID(id uint) (*googleCalendar.Token, error)
}

type googleCalendarTokenRepository struct {
	db *gorm.DB
}

func NewGCalendarRepository(db *gorm.DB) GoogleCalendarTokenRepository {
	return &googleCalendarTokenRepository{db: db}
}

func (s *googleCalendarTokenRepository) Store(token string, expirationDateTime time.Time, psyID uint) error {
	print("psyID: ", psyID)
	googleCalendarToken := googleCalendar.Token{
		PsyID:              psyID,
		Token:              token,
		ExpirationDateTime: expirationDateTime,
		TenantID:           1,
	}
	return s.db.Create(&googleCalendarToken).Error
}

func (s *googleCalendarTokenRepository) FindTokenByPsychologistID(id uint) (*googleCalendar.Token, error) {
	var token googleCalendar.Token
	result := s.db.Where("PsyID = ?", id).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func (s *googleCalendarTokenRepository) FindByPsyID(psyID uint) (*googleCalendar.TokenDTO, error) {
	var googleCalendarToken googleCalendar.Token
	err := s.db.Where("psy_id = ?", psyID).First(&googleCalendarToken).Error
	if err != nil {
		return nil, err
	}
	return tokenToDTO(&googleCalendarToken), nil
}

func (s *googleCalendarTokenRepository) Update(token string, psyID uint) error {
	return s.db.Model(&googleCalendar.Token{}).Where("psy_id = ?", psyID).Updates(googleCalendar.TokenDTO{
		Token:              token,
		ExpirationDateTime: utils.GetExpirationDate(),
	}).Error
}

// Delete It won't be used, just impl to complete the crud.
func (s *googleCalendarTokenRepository) Delete(psyID uint) error {
	return s.db.Where("psy_id = ?", psyID).Delete(&googleCalendar.TokenDTO{}).Error
}

func tokenToDTO(token *googleCalendar.Token) *googleCalendar.TokenDTO {
	return &googleCalendar.TokenDTO{
		ID:                 token.ID,
		PsyId:              token.PsyID,
		Token:              token.Token,
		ExpirationDateTime: token.ExpirationDateTime,
		TenantID:           token.TenantID,
	}
}
