package app

import (
	"tevian/domain"
	"tevian/domain/services"
	"tevian/services/config"
	face_cloud "tevian/services/face-cloud"
)

type ctx struct {
	services   domain.Services
	connection domain.Connection
}

type svs struct {
	config    services.Config
	faceCloud services.FaceCloud
}

func (s *svs) Config() services.Config {
	return s.config
}

func (s *svs) FaceCloud() services.FaceCloud {
	return s.faceCloud
}

func (c *ctx) Services() domain.Services {
	return c.services
}

func (c *ctx) Connection() domain.Connection {
	return c.connection
}
func (c *ctx) Make() domain.Context {
	return &ctx{
		services:   c.services,
		connection: c.connection,
	}
}
func InitCtx() *ctx {
	cfg, err := config.Make()
	if err != nil {
		panic(err)
	}

	faceCloud, err := face_cloud.GetAccessToken(cfg.CloudLogin(), cfg.CloudPassword())
	if err != nil {
		panic(err)
	}

	connection, err := InitDb(cfg)
	if err != nil {
		panic(err)
	}

	return &ctx{
		services: &svs{
			config:    cfg,
			faceCloud: faceCloud,
		},
		connection: connection,
	}
}
