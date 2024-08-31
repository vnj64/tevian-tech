package app

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"tevian/domain"
)

func MiddlewareAuthRequired(ctx domain.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "login and username required"})
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth header"})
		}

		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("cant decode payload: %v", err)})
		}

		credentials := string(payload)
		parts = strings.Split(credentials, ":")
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid payload type"})
		}

		username, password := parts[0], parts[1]
		if username != ctx.Services().Config().ServiceLogin() || password != ctx.Services().Config().ServicePassword() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not correct username or password"})
		}

		return c.Next()
	}
}
