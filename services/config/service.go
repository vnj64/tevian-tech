package config

import (
	"errors"
	"github.com/joho/godotenv"
	"tevian/domain/services"
)

type service struct {
	cloudLogin         string
	cloudPassword      string
	postgresqlHost     string
	postgresqlPort     string
	postgresqlUser     string
	postgresqlPassword string
	postgresqlDatabase string
	baseFaceCloudUrl   string
}

func Make() (services.Config, error) {
	envFile, _ := godotenv.Read(".env")
	cloudLogin := envFile["CLOUD_LOGIN"]
	cloudPassword := envFile["CLOUD_PASSWORD"]
	postgresqlHost := envFile["POSTGRESQL_HOST"]
	postgresqlPort := envFile["POSTGRESQL_PORT"]
	postgresqlUser := envFile["POSTGRESQL_USER"]
	postgresqlPassword := envFile["POSTGRESQL_PASSWORD"]
	postgresqlDatabase := envFile["POSTGRESQL_DATABASE"]
	baseFaceCloudUrl := envFile["BASE_FACE_CLOUD_URL"]

	if cloudLogin == "" || cloudPassword == "" || postgresqlHost == "" || postgresqlPort == "" || postgresqlUser == "" || postgresqlPassword == "" || postgresqlDatabase == "" || baseFaceCloudUrl == "" {
		return nil, errors.New("please check variables on .env")
	}

	return &service{
		cloudLogin:         cloudLogin,
		cloudPassword:      cloudPassword,
		postgresqlHost:     postgresqlHost,
		postgresqlPort:     postgresqlPort,
		postgresqlUser:     postgresqlUser,
		postgresqlPassword: postgresqlPassword,
		postgresqlDatabase: postgresqlDatabase,
		baseFaceCloudUrl:   baseFaceCloudUrl,
	}, nil
}

func (s *service) CloudLogin() string {
	return s.cloudLogin
}

func (s *service) CloudPassword() string {
	return s.cloudPassword
}

func (s *service) PostgresqlHost() string {
	return s.postgresqlHost
}

func (s *service) PostgresqlPort() string {
	return s.postgresqlPort
}

func (s *service) PostgresqlUser() string {
	return s.postgresqlUser
}

func (s *service) PostgresqlPassword() string {
	return s.postgresqlPassword
}

func (s *service) PostgresqlDatabase() string {
	return s.postgresqlDatabase
}

func (s *service) BaseFaceCloudUrl() string {
	return s.baseFaceCloudUrl
}
