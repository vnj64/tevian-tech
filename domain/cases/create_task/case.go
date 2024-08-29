package create_task

import (
	"fmt"
	"tevian/domain"
	"tevian/domain/models"
)

type Request struct {
	Task models.Task `json:"task"`
}
type Response struct {
	Id string `json:"id"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	// предполагается, что при добавлении задание, у него не может быть другого статуса
	r.Task.Status = models.StatusForming

	id, err := c.Connection().Task().Insert(r.Task)
	if err != nil {
		return nil, fmt.Errorf("error inserting task: %v", err)
	}

	return &Response{Id: id}, nil
}
