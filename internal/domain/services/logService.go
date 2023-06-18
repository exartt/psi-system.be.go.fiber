package services

import "github.com/sirupsen/logrus"

func logAndReturnError(action string, message string, err error) error {
	Logger.WithFields(logrus.Fields{"action": action}).Error(message, ": ", err)
	return err
}

func logInfo(action string, message string) {
	Logger.WithFields(logrus.Fields{"action": action}).Info(message)
}
