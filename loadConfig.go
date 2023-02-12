package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("[ERROR] failed to run function 'loadEnvFile': %v", err)
	}
}

func getEnvValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("[ERROR] failed to run function 'getEnvValue': the value of %s is empty", key)
	}
	return value
}
