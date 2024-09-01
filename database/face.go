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

func (f face) model() *models.Face {
	return &models.Face{
		Id:      f.Id,
		ImageId: f.ImageId,
		Bbox:    f.Bbox,
		Gender:  f.Gender,
		Age:     f.Age,
	}
}

func makeFace(f models.Face) *models.Face {
	return &models.Face{
		Id:      f.Id,
		ImageId: f.ImageId,
		Bbox:    f.Bbox,
		Gender:  f.Gender,
		Age:     f.Age,
	}
}

func (fr *faceRepository) Insert(face models.Face) (string, error) {
	tsk := makeFace(face)

	if err := fr.db.Create(tsk).Error; err != nil {
		return "", fmt.Errorf("unable to create face: %v", err)
	}

	return tsk.Id, nil
}

func (fr *faceRepository) WhereId(id string) (*models.Face, error) {
	var result face

	if err := fr.db.Where(models.Face{Id: id}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("unable to find face: %v", err)
	}

	return result.model(), nil
}

func (fr *faceRepository) WhereImageId(id string) ([]models.Face, error) {
	var results []face

	if err := fr.db.Where(models.Face{ImageId: id}).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("unable to find faces: %v", err)
	}

	out := make([]models.Face, len(results))

	for j, i := range results {
		out[j] = *i.model()
	}

	return out, nil
}

func (fr *faceRepository) Update(id string, updates map[string]interface{}) error {
	return fr.db.Model(&face{Id: id}).Updates(updates).Error
}

func (fr *faceRepository) Delete(id string) error {
	return fr.db.Delete(&face{Id: id}).Error
}
