package app

import (
	"tevian/database"
	"tevian/domain"
	"tevian/domain/services"
)

func InitDb(cfg services.Config) (domain.Connection, error) {
	return database.Make(cfg.PostgresqlUser(), cfg.PostgresqlPassword(), cfg.PostgresqlHost(), cfg.PostgresqlPort(), cfg.PostgresqlDatabase())
}
