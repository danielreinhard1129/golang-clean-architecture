package configs

import (
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/gofiber/fiber/v2"
)

func NewFiberConfiguration() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	}
}
