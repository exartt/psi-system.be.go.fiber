package config

import "os"

func ServerAddress() string {
	return os.Getenv("SERVER_ADDRESS")
}

func Secret() []byte {
	return []byte(os.Getenv("SECRET"))
}
