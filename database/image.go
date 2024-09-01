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

func (i image) model() *models.Image {
	return &models.Image{
		Id:           i.Id,
		TaskId:       i.TaskId,
		ImageName:    i.ImageName,
		ImageAddress: i.ImageAddress,
	}
}

func makeImage(i models.Image) *models.Image {
	return &models.Image{
		Id:           i.Id,
		TaskId:       i.TaskId,
		ImageName:    i.ImageName,
		ImageAddress: i.ImageAddress,
	}
}

func (ir *imageRepository) Insert(image models.Image) (string, error) {
	tsk := makeImage(image)

	if err := ir.db.Create(tsk).Error; err != nil {
		return "", fmt.Errorf("unable to create image: %v", err)
	}

	return tsk.Id, nil
}

func (ir *imageRepository) WhereId(id string) (*models.Image, error) {
	var result image

	if err := ir.db.Where(models.Image{Id: id}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("unable to find image: %v", err)
	}

	return result.model(), nil
}

func (ir *imageRepository) WhereTaskId(id string) ([]models.Image, error) {
	var results []image

	if err := ir.db.Where(models.Image{TaskId: id}).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("unable to find image: %v", err)
	}

	out := make([]models.Image, len(results))

	for j, i := range results {
		out[j] = *i.model()
	}

	return out, nil
}

func (ir *imageRepository) Update(id string, updates map[string]interface{}) error {
	return ir.db.Model(&image{Id: id}).Updates(updates).Error
}

func (ir *imageRepository) Delete(id string) error {
	return ir.db.Delete(&image{Id: id}).Error
}
