package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tevian/domain"
	"tevian/domain/repositories"
)

type connection struct {
	db *gorm.DB

	taskRepository repositories.Task
}

func makeConnection(db *gorm.DB) *connection {
	return &connection{
		db:             db,
		taskRepository: &taskRepository{db},
	}
}

func Make(user, password, host, port, database string, sslmode bool) (domain.Connection, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%v",
		user,
		password,
		host,
		port,
		database,
		true,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("unable to get database connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping DB: %v", err)
	}

	return makeConnection(db), nil
}

func (c connection) Task() repositories.Task {
	return c.taskRepository
}
