package util

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetDotenv(key string) string {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error load .env")
	}
	return os.Getenv(key)
}
