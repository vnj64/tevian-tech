package database

import (
	"fmt"
	"gorm.io/gorm"
	"tevian/domain/models"
)

type faceRepository struct {
	db *gorm.DB
}

type face struct {
	Id      string
	ImageId string
	Bbox    string
	Gender  string
	Age     float64
}

func (t face) model() *models.Face {
	return &models.Face{
		Id:      t.Id,
		ImageId: t.ImageId,
		Bbox:    t.Bbox,
		Gender:  t.Gender,
		Age:     t.Age,
	}
}

func makeFace(t models.Face) *models.Face {
	return &models.Face{
		Id:      t.Id,
		ImageId: t.ImageId,
		Bbox:    t.Bbox,
		Gender:  t.Gender,
		Age:     t.Age,
	}
}

func (tr *faceRepository) Insert(face models.Face) (string, error) {
	tsk := makeFace(face)

	if err := tr.db.Create(tsk).Error; err != nil {
		return "", fmt.Errorf("unable to create face: %v", err)
	}

	return tsk.Id, nil
}

func (tr *faceRepository) WhereId(id string) (*models.Face, error) {
	var result face

	if err := tr.db.Where(models.Face{Id: id}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("unable to find face: %v", err)
	}

	return result.model(), nil
}

func (tr *faceRepository) WhereImageId(id string) ([]models.Face, error) {
	var results []face

	if err := tr.db.Where(models.Face{ImageId: id}).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("unable to find faces: %v", err)
	}

	out := make([]models.Face, len(results))

	for j, i := range results {
		out[j] = *i.model()
	}

	return out, nil
}

func (tr *faceRepository) Update(id string, updates map[string]interface{}) error {
	return tr.db.Model(&face{Id: id}).Updates(updates).Error
}

func (tr *faceRepository) Delete(id string) error {
	return tr.db.Delete(&face{Id: id}).Error
}
