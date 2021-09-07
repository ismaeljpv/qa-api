package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

//Config method that return the Database conenction URI from the env file
//or an error when retrieving the data.
func GetConfig(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", errors.New("Failed to read configuration file.")
	}
	return os.Getenv(key), nil
}
