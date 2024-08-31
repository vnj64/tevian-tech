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
	_, err := c.Connection().Task().WhereId(r.Id)
	if err != nil {
		return nil, fmt.Errorf("task with id [%s] does not exist: %v", r.Id, err)
	}

	images, err := c.Connection().Image().WhereTaskId(r.Id)
	if err != nil {
		return nil, fmt.Errorf("error fetching images for task [%s]: %v", r.Id, err)
	}

	updates := make(map[string]interface{})
	updates["status"] = models.StatusProcessing
	if err := c.Connection().Task().Update(r.Id, updates); err != nil {
		return nil, fmt.Errorf("error updating task status to processing: %v", err)
	}

	go processTask(c, r.Id, images)

	return &Response{
		Message: "Task is successfully processing.",
	}, nil
}

func processTask(c domain.Context, taskId string, images []models.Image) {
	config := c.Services().Config()
	token, err := c.Services().FaceCloud().GetAccessToken(config.CloudLogin(), config.CloudPassword())
	if err != nil {
		updateTaskStatus(c, taskId, models.StatusFailed)
		fmt.Printf("error getting access token: %v\n", err)
		return
	}

	result, err := c.Services().FaceCloud().Detect(token, images)
	if err != nil {
		updateTaskStatus(c, taskId, models.StatusFailed)
		fmt.Printf("error during image detection: %v\n", err)
		return
	}

	imageIdMap := make(map[string]string)
	for _, image := range images {
		imageIdMap[image.ImageName] = image.Id
	}

	for _, imgData := range result.ImageData {
		imageId, exists := imageIdMap[imgData.Name]
		if !exists {
			fmt.Printf("no image ID found for image name: %s\n", imgData.Name)
			continue
		}

		for _, face := range imgData.Faces {
			faceModel := models.Face{
				Id:      uuid.New().String(),
				ImageId: imageId,
				Bbox:    fmt.Sprintf("%d,%d,%d,%d", face.BoundingBox.X, face.BoundingBox.Y, face.BoundingBox.Width, face.BoundingBox.Height),
				Gender:  face.Gender,
				Age:     face.Age,
			}

			if _, err := c.Connection().Face().Insert(faceModel); err != nil {
				updateTaskStatus(c, taskId, models.StatusFailed)
				fmt.Printf("error inserting face: %v\n", err)
				return
			}
		}
	}
	_, err = makeCalculateStatistics(c, taskId, images)
	if err != nil {
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

func makeCalculateStatistics(c domain.Context, taskId string, images []models.Image) (models.Statistics, error) {
	var stats models.Statistics
	var totalMaleAge, totalFemaleAge float64

	for _, image := range images {
		faces, err := c.Connection().Face().WhereImageId(image.Id)
		if err != nil {
			return models.Statistics{}, err
		}

		for _, face := range faces {
			stats.TotalFaces++
			if face.Gender == "male" {
				stats.TotalMales++
				totalMaleAge += face.Age
			} else if face.Gender == "female" {
				stats.TotalFemales++
				totalFemaleAge += face.Age
			}
		}
	}

	if stats.TotalMales > 0 {
		stats.AverageMaleAge = totalMaleAge / float64(stats.TotalMales)
	}

	if stats.TotalFemales > 0 {
		stats.AverageFemaleAge = totalFemaleAge / float64(stats.TotalFemales)
	}

	taskUpdates := map[string]interface{}{
		"all_faces_quantity": stats.TotalFaces,
		"male_quantity":      stats.TotalMales,
		"female_quantity":    stats.TotalFemales,
		"average_male_age":   stats.AverageMaleAge,
		"average_female_age": stats.AverageFemaleAge,
	}

	if err := c.Connection().Task().Update(taskId, taskUpdates); err != nil {
		updateTaskStatus(c, taskId, models.StatusFailed)
		fmt.Printf("error updating task statistics: %v", err)
	}

	return stats, nil
}
