package utils

import (
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
	case "CASHFLOW":
		return enums.CASHFLOW, nil
	default:
		return enums.RECEIVABLE, nil
	}
}
