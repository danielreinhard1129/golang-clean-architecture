package handler

import (
	"strconv"

	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/middleware"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/request"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/usecase"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/pagination"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase.UserUsecase
}

func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: *usecase}
}

func (h *UserHandler) Route(app *fiber.App) {
	app.Get("/v1/users", middleware.JWTProtected(), middleware.RequireRoles("USER"), h.FindAll)
	app.Get("/v1/users/:id", h.FindById)
	app.Post("/v1/users", h.Create)
	app.Patch("/v1/users/:id", h.Update)
	app.Delete("/v1/users/:id", h.Delete)
}

func (h *UserHandler) FindAll(c *fiber.Ctx) error {
	qp := request.ParseAndValidate(c)

	result, total := h.UserUsecase.FindAll(qp.Search, qp.OrderBy, qp.Sort, qp.Page, qp.Limit)

	return c.Status(fiber.StatusOK).JSON(pagination.Response[entities.User]{
		Data: result,
		Meta: pagination.Meta{
			Page:  qp.Page,
			Limit: qp.Limit,
			Total: total,
		},
	})
}

func (h *UserHandler) FindById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)

	result, err := h.UserUsecase.FindById(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var reqBody request.UserCreateRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	result, err := h.UserUsecase.Create(&reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var reqBody request.UserUpdateRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	id, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	result, err := h.UserUsecase.Update(id, &reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)

}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)

	err = h.UserUsecase.Delete(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}
