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

func GetFifthDayOfCurrentMonth() string {
	firstDayOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	fifthDayOfMonth := firstDayOfMonth.AddDate(0, 0, 4)
	return fifthDayOfMonth.Format("2006-01-02")
}
