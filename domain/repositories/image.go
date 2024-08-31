package repositories

import "tevian/domain/models"

type Image interface {
	Insert(image models.Image) (string, error)
	WhereId(id string) (*models.Image, error)
	WhereTaskId(id string) ([]models.Image, error)
	Delete(id string) error
	Update(id string, updates map[string]interface{}) error
}
