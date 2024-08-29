package start_task

import (
	"fmt"
	"tevian/domain"
	"tevian/domain/models"
)

type Request struct {
	Id string `json:"id"`
}

type Response struct {
	Result []*models.ResultData `json:"result"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	config := c.Services().Config()

	token, err := c.Services().FaceCloud().GetAccessToken(config.CloudLogin(), config.CloudPassword())
	if err != nil {
		return nil, fmt.Errorf("error getting access token: %v", err)
	}

	task, err := c.Connection().Task().WhereId(r.Id)
	if err != nil {
		return nil, fmt.Errorf("task with id [%s] does not exist: %v", r.Id, err)
	}

	var imageAddresses []string
	imageAddresses = append(imageAddresses, task.ImageAddress)

	results, err := c.Services().FaceCloud().Detect(token, imageAddresses)
	if err != nil {
		return nil, fmt.Errorf("error detecting faces: %v", err)
	}

	return &Response{
		Result: results,
	}, nil
}
