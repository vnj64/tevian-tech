package repositories

import "tevian/domain/models"

type Face interface {
	Insert(face models.Face) (string, error)
	WhereId(id string) (*models.Face, error)
	WhereImageId(id string) ([]models.Face, error)
	Delete(id string) error
	Update(id string, updates map[string]interface{}) error
}
