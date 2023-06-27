package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	AppPort              int
	DbHost               string
	DbPort               int
	DbUserName           string
	DbPassword           string
	DbDatabase           string
	DbRunMigrations      bool
	DbMigrationsLocation string
}

func GetConfig(envFilePath string) AppConfig {
	log.Infof("envFilePath: %s", envFilePath)
	if _, err := os.Stat(envFilePath); errors.Is(err, fs.ErrNotExist) {
		log.Infof("%s file doesn't exist", envFilePath)
	} else {
		err := godotenv.Load(envFilePath)
		if err != nil {
			log.Warningf("Couldn't load environment variables from .env file: %s", envFilePath)
		}
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	AppPort, _ := strconv.Atoi(port)
	DbHost := os.Getenv("APP_DB_HOST")
	DbPort, _ := strconv.Atoi(os.Getenv("APP_DB_PORT"))
	DbUserName := os.Getenv("APP_DB_USERNAME")
	DbPassword := os.Getenv("APP_DB_PASSWORD")
	DbDatabase := os.Getenv("APP_DB_NAME")
	DbRunMigrations, _ := strconv.ParseBool(os.Getenv("APP_DB_RUN_MIGRATIONS"))
	DbMigrationsLocation := os.Getenv("APP_DB_MIGRATIONS_LOCATION")
	return AppConfig{
		AppPort:              AppPort,
		DbHost:               DbHost,
		DbPort:               DbPort,
		DbUserName:           DbUserName,
		DbPassword:           DbPassword,
		DbDatabase:           DbDatabase,
		DbRunMigrations:      DbRunMigrations,
		DbMigrationsLocation: DbMigrationsLocation,
	}
}
