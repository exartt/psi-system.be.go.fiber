package utils

import (
	"errors"
	"psi-system.be.go.fiber/internal/domain/enums"
	"time"
)

func GetExpirationDate() time.Time {
	currentTime := time.Now()
	return currentTime.Add(30 * time.Minute)
}

func CastToTransactioTypeEnum(typeParam string) (enums.TransactionType, error) {
	switch typeParam {
	case "PAYABLE":
		return enums.PAYABLE, nil
	case "RECEIVABLE":
		return enums.RECEIVABLE, nil
	case "CASHFLOW":
		return enums.CASHFLOW, nil
	default:
		return 500, errors.New("Não foi possível identificar o tipo de transação, por favor atualize a página e tente novamente")
	}
}
