package common

import (
	"os"
)

type config struct {
	DatabaseHost     string
	DatabasePort     string
	DatabasePassword string
	DatabaseUser     string
	DatabaseName     string
	Environment      string
	GitVersion       string
	WebAddress       string
	AwsBucket        string
	KmsARN           string
}

// AppConfig is the application configuration
var AppConfig config

func loadConfig() {
	AppConfig.DatabaseHost = os.Getenv("DATABASE_HOST")
	AppConfig.DatabasePort = os.Getenv("DATABASE_PORT")
	AppConfig.DatabasePassword = os.Getenv("DATABASE_PASSWORD")
	AppConfig.DatabaseUser = os.Getenv("DATABASE_USER")
	AppConfig.DatabaseName = os.Getenv("DATABASE_NAME")
	AppConfig.Environment = os.Getenv("ENVIRONMENT")
	AppConfig.GitVersion = os.Getenv("GIT_VERSION")
	AppConfig.WebAddress = os.Getenv("WEB_ADDRESS")
	AppConfig.AwsBucket = os.Getenv("AWS_BUCKET")
	AppConfig.KmsARN = os.Getenv("KMS_ARN")
}
