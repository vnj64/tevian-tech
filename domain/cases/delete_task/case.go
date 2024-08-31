package delete_task

import (
	"errors"
	"fmt"
	"os"
	"tevian/domain"
	"tevian/domain/models"
)

type Request struct {
	Id string `json:"id"`
}

func Run(c domain.Context, r Request) error {
	task, err := c.Connection().Task().WhereId(r.Id)
	if err != nil {
		return fmt.Errorf("task with id [%s] not found", r.Id)
	}

	if task.Status == models.StatusProcessing {
		return errors.New("can't delete task which is still processing")
	}

	images, err := c.Connection().Image().WhereTaskId(task.Id)
	if err != nil {
		return fmt.Errorf("error fetching images for task [%s]: %v", r.Id, err)
	}

	for _, image := range images {
		if image.ImageAddress != "" {
			if err := os.Remove(image.ImageAddress); err != nil {
				return fmt.Errorf("cannot delete image file [%s]: %v", image.ImageAddress, err)
			}
		}
	}

	if err := c.Connection().Task().Delete(r.Id); err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	return nil
}
