package utils

import (
	"fmt"
	"os"
)

func GenerateDSN() string {
	var url string
	if os.Getenv("DSN") != "" {
		url = os.Getenv("DSN")
	} else {
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL_MODE"),
		)
	}

	return url
}
