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

	taskRepository  repositories.Task
	faceRepository  repositories.Face
	imageRepository repositories.Image
}

func makeConnection(db *gorm.DB) *connection {
	return &connection{
		db:              db,
		taskRepository:  &taskRepository{db},
		faceRepository:  &faceRepository{db},
		imageRepository: &imageRepository{db},
	}
}

func Make(user, password, host, port, database string) (domain.Connection, error) {
	// postgresql://tevian:tevian@db:5432/postgres
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		database,
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

func (c connection) Face() repositories.Face {
	return c.faceRepository
}

func (c connection) Image() repositories.Image {
	return c.imageRepository
}
