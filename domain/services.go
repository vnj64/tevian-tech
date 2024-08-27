package domain

import "tevian/domain/services"

type Services interface {
	FaceCloud() services.FaceCloud
}
