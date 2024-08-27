package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"runtime"
	"strings"
)

type HttpServer struct {
	app *fiber.App
}

type Server interface {
	Start()
}

func NewHttpServer() Server {
	app := fiber.New(fiber.Config{
		BodyLimit:         1024 * 1024 * 50,
		AppName:           "TevianApp",
		StreamRequestBody: true,
	})

	var methods = []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions}
	var headers = []string{fiber.HeaderAccept, fiber.HeaderAuthorization, fiber.HeaderContentType,
		fiber.HeaderContentLength, fiber.HeaderAcceptEncoding, "X-CSRF-Token"}

	corsConfig := cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     strings.Join(methods, ","),
		AllowHeaders:     strings.Join(headers, ", "),
		AllowCredentials: true,
		MaxAge:           300,
	})

	app.Use(corsConfig)
	app.Use(recover.New())

	return &HttpServer{app: app}
}

func (s *HttpServer) Start() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	domainContext := InitCtx().Make()

	s.app.Use("", func(ctx *fiber.Ctx) error {
		ctx.Locals("context", domainContext)
		return ctx.Next()
	})

	err := s.app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
