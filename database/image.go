package database

import (
	"fmt"
	"gorm.io/gorm"
	"tevian/domain/models"
)

type imageRepository struct {
	db *gorm.DB
}

type image struct {
	Id           string
	TaskId       string
	ImageName    string
	ImageAddress string
}

func (t image) model() *models.Image {
	return &models.Image{
		Id:           t.Id,
		TaskId:       t.TaskId,
		ImageName:    t.ImageName,
		ImageAddress: t.ImageAddress,
	}
}

func makeImage(t models.Image) *models.Image {
	return &models.Image{
		Id:           t.Id,
		TaskId:       t.TaskId,
		ImageName:    t.ImageName,
		ImageAddress: t.ImageAddress,
	}
}

func (tr *imageRepository) Insert(image models.Image) (string, error) {
	tsk := makeImage(image)

	if err := tr.db.Create(tsk).Error; err != nil {
		return "", fmt.Errorf("unable to create image: %v", err)
	}

	return tsk.Id, nil
}

func (tr *imageRepository) WhereId(id string) (*models.Image, error) {
	var result image

	if err := tr.db.Where(models.Image{Id: id}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("unable to find image: %v", err)
	}

	return result.model(), nil
}

func (tr *imageRepository) WhereTaskId(id string) ([]models.Image, error) {
	var results []image

	if err := tr.db.Where(models.Image{TaskId: id}).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("unable to find image: %v", err)
	}

	out := make([]models.Image, len(results))

	for j, i := range results {
		out[j] = *i.model()
	}

	return out, nil
}

func (tr *imageRepository) Update(id string, updates map[string]interface{}) error {
	return tr.db.Model(&image{Id: id}).Updates(updates).Error
}

func (tr *imageRepository) Delete(id string) error {
	return tr.db.Delete(&image{Id: id}).Error
}
