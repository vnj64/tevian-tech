package upload_task_image

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"tevian/domain"
	"tevian/domain/models"
)

type Request struct {
	Id        string `json:"id"`
	ImageName string `json:"imageName"`
	Body      []byte `json:"body"`
}

func Run(c domain.Context, r Request) error {
	task, err := c.Connection().Task().WhereId(r.Id)
	if err != nil {
		return fmt.Errorf("task with id [%s] does not exist due [%v]", r.Id, err)
	}

	if task.Status == models.StatusProcessing || task.Status == models.StatusCompleted {
		return fmt.Errorf("unavailable to upload image to this task: status is %s", task.Status)
	}

	r.ImageName = uuid.New().String() + ".jpg"
	file, err := os.Create(fmt.Sprintf("images/%s", r.ImageName))
	if err != nil {
		return fmt.Errorf("failed to create file due [%v]", err)
	}
	defer file.Close()

	reader := bytes.NewReader(r.Body)

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write to file due [%v]", err)
	}

	updates := make(map[string]interface{})

	updates["image_name"] = r.ImageName
	updates["image_address"] = fmt.Sprintf("images/%s", r.ImageName)

	if len(updates) == 0 {
		return fmt.Errorf("nothing to update")
	}

	img := models.Image{
		Id:           uuid.New().String(),
		TaskId:       r.Id,
		ImageName:    r.ImageName + ".jpg",
		ImageAddress: fmt.Sprintf("images/%s", r.ImageName),
	}

	_, err = c.Connection().Image().Insert(img)
	if err != nil {
		return fmt.Errorf("failed to insert image due [%v]", err)
	}

	return nil
}
