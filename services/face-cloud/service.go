package face_cloud

import (
	"tevian/domain/services"
)

type service struct {
	cfg services.Config
}

func Make(cfg services.Config) services.FaceCloud {
	return &service{
		cfg: cfg,
	}
}
