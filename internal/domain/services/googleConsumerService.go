package services

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/googleCalendar"
	"psi-system.be.go.fiber/internal/domain/utils"
	"psi-system.be.go.fiber/internal/repositories"
	"time"
)

type GoogleConsumerService interface {
	RequestGoogleAuth() (string, error)
	RequestGoogleAuthAuthorized() (string, error)
	Store(token string, expirationDateTime time.Time, psyID uint) error
	FindByPsyID(psyID uint) (*googleCalendar.TokenDTO, error)
	Update(token string, psyID uint) error
	Delete(psyID uint, tenantID uint) error
	ExchangeCodeAndStoreToken(code string, psychologistID uint) (*googleCalendar.TokenDTO, error)
	FindTokenByPsychologistID(id uint) (*googleCalendar.Token, error)
}

func NewGCalendarService(repo repositories.GoogleCalendarTokenRepository) GoogleConsumerService {
	return &googleConsumerService{repo: repo}
}

type googleConsumerService struct {
	repo repositories.GoogleCalendarTokenRepository
}

func (service *googleConsumerService) Store(token string, expirationDateTime time.Time, psychologistID uint) error {
	_, err := service.repo.FindByPsyID(psychologistID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return service.repo.Store(token, expirationDateTime, psychologistID)
		}
		return err
	}

	return service.repo.Update(token, psychologistID)
}

func (service *googleConsumerService) FindTokenByPsychologistID(id uint) (*googleCalendar.Token, error) {
	return service.repo.FindTokenByPsychologistID(id)
}

func (service *googleConsumerService) FindByPsyID(psychologistID uint) (*googleCalendar.TokenDTO, error) {
	tokenDTO, err := service.repo.FindByPsyID(psychologistID)
	currentTime := time.Now()

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newToken, err := service.RequestGoogleAuth()
			expirationTime := utils.GetExpirationDate()
			if err != nil {
				return nil, err
			}
			err = service.Store(newToken, expirationTime, psychologistID)
			if err != nil {
				return nil, err
			}
			return &googleCalendar.TokenDTO{
				PsyId:              psychologistID,
				Token:              newToken,
				ExpirationDateTime: expirationTime,
				TenantID:           0,
			}, nil
		}
		return nil, err
	}

	if tokenDTO.ExpirationDateTime.Before(currentTime) {
		newToken, err := service.RequestGoogleAuthAuthorized()
		expirationTime := utils.GetExpirationDate()
		if err != nil {
			return nil, err
		}
		err = service.Store(newToken, expirationTime, psychologistID)
		if err != nil {
			return nil, err
		}
		return &googleCalendar.TokenDTO{
			PsyId:              psychologistID,
			Token:              newToken,
			ExpirationDateTime: expirationTime,
			TenantID:           0,
		}, nil
	}

	return tokenDTO, nil
}

func (service *googleConsumerService) Update(token string, psychologistID uint) error {
	return service.repo.Update(token, psychologistID)
}

func (service *googleConsumerService) Delete(psychologistID uint, tenantID uint) error {
	return service.repo.Delete(psychologistID)
}

func (service *googleConsumerService) RequestGoogleAuth() (string, error) {
	resp, err := http.Get(enums.AUTH.String(""))
	if err != nil {
		logAndReturnError("GoogleAuth", "Erro ao fazer requisição para o serviço2: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logAndReturnError("GoogleAuth", "Erro ao ler a resposta do serviço2: %v", err)
		return "", err
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		logAndReturnError("GoogleAuth", "Ocorreu um erro inesperado com a API CALENDAR", err)
		return "", err
	}
	authURL := result["auth_url"]

	return authURL, nil
}

func (service *googleConsumerService) RequestGoogleAuthAuthorized() (string, error) {
	resp, err := http.Get(enums.AUTH.String(""))
	if err != nil {
		logAndReturnError("GoogleAuth", "Erro ao fazer requisição para o serviço: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logAndReturnError("GoogleAuth", "Erro ao ler a resposta do serviço: %v", err)
		return "", err
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		logAndReturnError("GoogleAuth", "Ocorreu um erro inesperado com a API CALENDAR", err)
		return "", err
	}

	authURL := result["auth_url"]

	respToken, err := http.Get(authURL)
	if err != nil {
		logAndReturnError("GoogleAuth", "Erro ao fazer requisição para a URL do token: %v", err)
		return "", err
	}
	defer respToken.Body.Close()

	tokenBody, err := ioutil.ReadAll(respToken.Body)
	if err != nil {
		logAndReturnError("GoogleAuth", "Erro ao ler a resposta da URL do token: %v", err)
		return "", err
	}

	var tokenResult map[string]string
	err = json.Unmarshal(tokenBody, &tokenResult)
	if err != nil {
		logAndReturnError("GoogleAuth", "Erro ao deserializar o JSON do token: %v", err)
		return "", err
	}

	token := tokenResult["access_token"]

	return token, nil
}

func (service *googleConsumerService) ExchangeCodeAndStoreToken(code string, psychologistID uint) (*googleCalendar.TokenDTO, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:3021/calendar/v1/callback?code=%s", code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	accessToken := result["access_token"]

	expirationTime := utils.GetExpirationDate()
	err = service.Store(accessToken, expirationTime, psychologistID)
	if err != nil {
		return nil, err
	}

	return &googleCalendar.TokenDTO{
		PsyId:              psychologistID,
		Token:              accessToken,
		ExpirationDateTime: expirationTime,
		TenantID:           0,
	}, nil
}
