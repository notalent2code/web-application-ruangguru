package config

import "os"

var (
	BaseURL = os.Getenv("BASE_URL")
	PORT    = os.Getenv("PORT")
)

func SetUrl(url string) string {
	if BaseURL == "" {
		BaseURL = "http://localhost:8081"
	}

	return BaseURL + url
}

func SetPort() string {
	if PORT == "" {
		PORT = "8081"
	}

	return PORT
}