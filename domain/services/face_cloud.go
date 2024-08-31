package services

import (
	"tevian/domain/models"
)

type FaceCloud interface {
	Detect(token string, tasks *models.Task) (models.DetectResult, error)
	GetAccessToken(cloudLogin, cloudPassword string) (string, error)
}
