package domain

import "tevian/domain/services"

type Services interface {
	Config() services.Config
	FaceCloud() services.FaceCloud
}
