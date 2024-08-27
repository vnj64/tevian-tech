package services

type Config interface {
	CloudLogin() string
	CloudPassword() string
	PostgresqlHost() string
	PostgresqlPort() string
	PostgresqlUser() string
	PostgresqlPassword() string
	PostgresqlDatabase() string
	BaseFaceCloudUrl() string
}
