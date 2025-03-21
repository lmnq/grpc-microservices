package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("failed to parse port number: %v", err)
	}
	return port
}

func getEnvironmentValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s not set", key)
	}
	return value
}
