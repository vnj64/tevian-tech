package services

import (
	"tevian/domain/models"
)

type FaceCloud interface {
	Detect(token string, imageAddresses []string) ([]*models.ResultData, error)
	GetAccessToken(cloudLogin, cloudPassword string) (string, error)
}
