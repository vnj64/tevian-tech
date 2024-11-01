package database

import (
	"fmt"
	"gorm.io/gorm"
	"tevian/domain/models"
)

type taskRepository struct {
	db *gorm.DB
}

type task struct {
	Id               string
	Status           models.TaskStatus
	AllFacesQuantity *int
	MaleQuantity     *int
	FemaleQuantity   *int
	AverageMaleAge   *float64
	AverageFemaleAge *float64
}

func (t task) model() *models.Task {
	return &models.Task{
		Id:               t.Id,
		Status:           t.Status,
		AllFacesQuantity: t.AllFacesQuantity,
		MaleQuantity:     t.MaleQuantity,
		FemaleQuantity:   t.FemaleQuantity,
		AverageMaleAge:   t.AverageMaleAge,
		AverageFemaleAge: t.AverageFemaleAge,
	}
}

func makeTask(t models.Task) *models.Task {
	return &models.Task{
		Id:               t.Id,
		Status:           t.Status,
		AllFacesQuantity: t.AllFacesQuantity,
		MaleQuantity:     t.MaleQuantity,
		FemaleQuantity:   t.FemaleQuantity,
		AverageMaleAge:   t.AverageMaleAge,
		AverageFemaleAge: t.AverageFemaleAge,
	}
}

func (tr *taskRepository) Insert(task models.Task) (string, error) {
	tsk := makeTask(task)

	if err := tr.db.Create(tsk).Error; err != nil {
		return "", fmt.Errorf("unable to create task: %v", err)
	}

	return tsk.Id, nil
}

func (tr *taskRepository) WhereId(id string) (*models.Task, error) {
	var result task

	if err := tr.db.Where(models.Task{Id: id}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("unable to find task: %v", err)
	}

	return result.model(), nil
}

func (tr *taskRepository) Update(id string, updates map[string]interface{}) error {
	return tr.db.Model(&task{Id: id}).Updates(updates).Error
}

func (tr *taskRepository) Delete(id string) error {
	return tr.db.Delete(&task{Id: id}).Error
}
