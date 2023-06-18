package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/enums"
)

type GoogleConsumerService interface {
	RequestGoogleAuth() (string, error)
}

type googleConsumerService struct{}

func NewGoogleConsumerService() GoogleConsumerService {
	return &googleConsumerService{}
}

func (g *googleConsumerService) RequestGoogleAuth() (string, error) {
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
	json.Unmarshal(body, &result)
	authURL := result["auth_url"]

	return authURL, nil
}
