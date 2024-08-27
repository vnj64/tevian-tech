package config

import (
	"errors"
	"os"
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
	cloudLogin := os.Getenv("CLOUD_LOGIN")
	cloudPassword := os.Getenv("CLOUD_PASSWORD")
	postgresqlHost := os.Getenv("POSTGRESQL_HOST")
	postgresqlPort := os.Getenv("POSTGRESQL_PORT")
	postgresqlUser := os.Getenv("POSTGRESQL_USER")
	postgresqlPassword := os.Getenv("POSTGRESQL_PASSWORD")
	postgresqlDatabase := os.Getenv("POSTGRESQL_DATABASE")
	baseFaceCloudUrl := os.Getenv("BASE_FACE_CLOUD_URL")

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
