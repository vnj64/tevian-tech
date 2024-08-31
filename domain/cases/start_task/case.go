package start_task

import (
	"fmt"
	"github.com/google/uuid"
	"tevian/domain"
	"tevian/domain/models"
)

type Request struct {
	Id string `json:"id"`
}

type Response struct {
	Message string `json:"message"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	task, err := c.Connection().Task().WhereId(r.Id)
	if err != nil {
		return nil, fmt.Errorf("task with id [%s] does not exist: %v", r.Id, err)
	}

	if task.ImageAddress == nil {
		return nil, fmt.Errorf("image address for task with id [%s] is not set", r.Id)
	}

	updates := make(map[string]interface{})
	updates["status"] = models.StatusProcessing
	if err := c.Connection().Task().Update(r.Id, updates); err != nil {
		return nil, fmt.Errorf("error updating task status to processing: %v", err)
	}

	go processTask(c, r.Id, task)

	return &Response{
		Message: "Task is successfully processing.",
	}, nil
}

func processTask(c domain.Context, taskId string, task *models.Task) {
	config := c.Services().Config()
	token, err := c.Services().FaceCloud().GetAccessToken(config.CloudLogin(), config.CloudPassword())
	if err != nil {
		updateTaskStatus(c, taskId, models.StatusFailed)
		fmt.Printf("error getting access token: %v\n", err)
		return
	}

	result, err := c.Services().FaceCloud().Detect(token, task)
	if err != nil {
		updateTaskStatus(c, taskId, models.StatusFailed)
		fmt.Printf("error during image detection: %v\n", err)
		return
	}

	if err := saveTaskResults(c, taskId, result); err != nil {
		updateTaskStatus(c, taskId, models.StatusFailed)
		fmt.Printf("error saving task results: %v\n", err)
		return
	}

	updateTaskStatus(c, taskId, models.StatusCompleted)
}

func updateTaskStatus(c domain.Context, taskId string, status models.TaskStatus) {
	updates := map[string]interface{}{
		"status": status,
	}
	if err := c.Connection().Task().Update(taskId, updates); err != nil {
		fmt.Printf("error updating task status: %v\n", err)
	}
}

func saveTaskResults(c domain.Context, taskId string, result models.DetectResult) error {
	for _, imgData := range result.ImageData {
		image := models.Image{
			Id:     uuid.New().String(),
			TaskId: taskId,
			Name:   imgData.Name,
		}

		imageId, err := c.Connection().Image().Insert(image)
		if err != nil {
			return fmt.Errorf("error inserting image: %v", err)
		}

		for _, face := range imgData.Faces {
			faceModel := models.Face{
				Id:      uuid.New().String(),
				ImageId: imageId,
				Bbox:    fmt.Sprintf("%d,%d,%d,%d", face.BoundingBox.X, face.BoundingBox.Y, face.BoundingBox.Width, face.BoundingBox.Height),
				Gender:  face.Gender,
				Age:     int(face.Age),
			}

			if _, err := c.Connection().Face().Insert(faceModel); err != nil {
				return fmt.Errorf("error inserting face: %v", err)
			}
		}
	}

	taskUpdates := map[string]interface{}{
		"all_faces_quantity": result.Statistics.TotalFaces,
		"male_quantity":      result.Statistics.TotalMales,
		"female_quantity":    result.Statistics.TotalFemales,
		"average_male_age":   result.Statistics.AverageMaleAge,
		"average_female_age": result.Statistics.AverageFemaleAge,
	}

	if err := c.Connection().Task().Update(taskId, taskUpdates); err != nil {
		return fmt.Errorf("error updating task statistics: %v", err)
	}

	return nil
}
