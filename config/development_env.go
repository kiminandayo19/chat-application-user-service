package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvDevType struct {
	APP_ENV     string
	APP_HOST    string
	APP_PORT    string
	APP_ENTRY   string
	DB_HOST     string
	DB_NAME     string
	DB_PASS     string
	DB_PORT     int
	DB_USERNAME string
	APIKEY      string
	JWT_SECRET  string
	REDIS_ADDR  string
	REDIS_PASS  string
	REDIS_DB    int
}

func NewEnvDev() (EnvDevType, error) {
	var env EnvDevType
	var err error = nil

	err = godotenv.Load()
	if err != nil {
		log.Panic("[env] - Failed to load development env", err)
	}

	env.APP_ENV = os.Getenv("APP_ENV")
	env.APP_HOST = os.Getenv("APP_HOST")
	env.APP_PORT = os.Getenv("APP_PORT")
	env.APP_ENTRY = os.Getenv("APP_ENTRY")

	env.DB_HOST = os.Getenv("DB_HOST")
	env.DB_NAME = os.Getenv("DB_NAME")
	env.DB_PASS = os.Getenv("DB_PASS")
	env.DB_USERNAME = os.Getenv("DB_USERNAME")
	env.DB_PORT, err = strconv.Atoi(os.Getenv("DB_PORT"))

	env.APIKEY = os.Getenv("API_KEY")
	env.JWT_SECRET = os.Getenv("JWT_SECRET")

	env.REDIS_ADDR = os.Getenv("REDIS_ADDR")
	env.REDIS_PASS = os.Getenv("REDIS_PASS")
	env.REDIS_DB, err = strconv.Atoi(os.Getenv("REDIS_DB"))

	log.Print("[env] - Success to load development env")
	return env, err
}
