package handler

import (
	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/request"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/response"
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

func (h *AuthHandler) Route(app *fiber.App) {
	app.Post("/v1/auth/login", h.Login)
	app.Post("/v1/auth/register", h.Register)
	app.Post("/v1/auth/verify", h.VerifyAccount)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var reqBody request.AuthLoginRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	result, err := h.AuthUsecase.Login(c.UserContext(), &reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.AuthLoginResponse{
		User:  result.User,
		Token: result.Token,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var reqBody request.AuthRegisterRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	err = h.AuthUsecase.Register(c.UserContext(), &reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.AuthRegisterResponse{
		Message: "Register success",
	})
}

func (h *AuthHandler) VerifyAccount(c *fiber.Ctx) error {
	var reqBody request.AuthVerifyAccountRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	err = h.AuthUsecase.VerifyAccount(c.UserContext(), &reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.AuthVerifyAccountResponse{
		Message: "Verify account success",
	})
}
