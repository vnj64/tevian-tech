package services

import (
	"tevian/domain/models"
)

type FaceCloud interface {
	Detect(token string, images []models.Image) (models.DetectResult, error)
	GetAccessToken(cloudLogin, cloudPassword string) (string, error)
}
