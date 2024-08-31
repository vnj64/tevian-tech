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
		return errors.New("cant delete task, which still processing")
	}

	if task.ImageAddress != nil {
		if err := os.Remove(*task.ImageAddress); err != nil {
			return fmt.Errorf("cannot delete this image task: %v", err)
		}
	}

	if err := c.Connection().Task().Delete(r.Id); err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	return nil
}
