package handler

import (
	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/request"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/usecase"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	usecase.AuthUsecase
}

func NewAuthHandler(usecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{AuthUsecase: *usecase}
}

func (h AuthHandler) Route(app *fiber.App) {
	app.Post("/v1/auth/login", h.Login)
}

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	var reqBody request.AuthLoginRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	result, err := handler.AuthUsecase.Login(&reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
