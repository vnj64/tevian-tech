package create_task

import "tevian/domain/models"

type Request struct {
	Task models.Task `json:"task"`
}
type Response struct {
	Id string `json:"id"`
}
