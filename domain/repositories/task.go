package repositories

import "tevian/domain/models"

type Task interface {
	Insert(task models.Task) (string, error)
	WhereId(id string) (*models.Task, error)
	Delete(id string) error
	Update(id string, updates map[string]interface{}) error
}
