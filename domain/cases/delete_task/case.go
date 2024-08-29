package delete_task

import (
	"errors"
	"fmt"
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

	if err := c.Connection().Task().Delete(r.Id); err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	return nil
}
