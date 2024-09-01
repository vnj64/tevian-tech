package get_task

import (
	"fmt"
	"tevian/domain"
	"tevian/domain/models"
)

type Request struct {
	Id string `json:"id"`
}

type Response struct {
	TaskId     string                   `json:"taskId"`
	TaskStatus models.TaskStatus        `json:"taskStatus"`
	Faces      map[string][]models.Face `json:"faces"`
	Statistics StatsAdditional          `json:"statistics"`
}

type StatsAdditional struct {
	AllFacesQuantity int     `json:"allFacesQuantity"`
	MaleQuantity     int     `json:"maleQuantity"`
	FemaleQuantity   int     `json:"femaleQuantity"`
	AverageMaleAge   float64 `json:"averageMaleAge"`
	AverageFemaleAge float64 `json:"averageFemaleAge"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	task, err := c.Connection().Task().WhereId(r.Id)
	if err != nil {
		return nil, fmt.Errorf("task with id [%s] does not exist: %v", r.Id, err)
	}

	images, err := c.Connection().Image().WhereTaskId(task.Id)
	if err != nil {
		return nil, fmt.Errorf("error fetching images for task with id [%s]: %v", r.Id, err)
	}

	facesByImage := make(map[string][]models.Face)
	for _, image := range images {
		imageFaces, err := c.Connection().Face().WhereImageId(image.Id)
		if err != nil {
			return nil, fmt.Errorf("error fetching faces for image with id [%s]: %v", image.Id, err)
		}
		facesByImage[image.ImageName] = imageFaces
	}

	return &Response{
		TaskId:     task.Id,
		TaskStatus: task.Status,
		Faces:      facesByImage,
		Statistics: StatsAdditional{
			AllFacesQuantity: *task.AllFacesQuantity,
			MaleQuantity:     *task.MaleQuantity,
			FemaleQuantity:   *task.FemaleQuantity,
			AverageMaleAge:   *task.AverageMaleAge,
			AverageFemaleAge: *task.AverageFemaleAge,
		},
	}, nil
}
