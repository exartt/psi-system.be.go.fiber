package utils

import "time"

func GetExpirationDate() time.Time {
	currentTime := time.Now()
	return currentTime.Add(30 * time.Minute)
}
