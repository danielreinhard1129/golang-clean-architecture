package handler

import (
	"strconv"

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

func (h UserHandler) Route(app *fiber.App) {
	app.Get("/v1/users", h.FindAll)
	app.Get("/v1/users/:id", h.FindById)
	app.Post("/v1/users", h.Create)
	app.Patch("/v1/users/:id", h.Update)
	app.Delete("/v1/users/:id", h.Delete)
}

func (handler UserHandler) FindAll(c *fiber.Ctx) error {
	qp := request.ParseAndValidate(c)

	result, total := handler.UserUsecase.FindAll(qp.Search, qp.OrderBy, qp.Sort, qp.Page, qp.Limit)

	return c.Status(fiber.StatusOK).JSON(pagination.Response[entities.User]{
		Data: result,
		Meta: pagination.Meta{
			Page:  qp.Page,
			Limit: qp.Limit,
			Total: total,
		},
	})
}

func (handler UserHandler) FindById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)

	result, err := handler.UserUsecase.FindById(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler UserHandler) Create(c *fiber.Ctx) error {
	var reqBody request.UserCreateRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	result, err := handler.UserUsecase.Create(&reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler UserHandler) Update(c *fiber.Ctx) error {
	var reqBody request.UserUpdateRequest
	err := c.BodyParser(&reqBody)
	exception.PanicLogging(err)

	id, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)

	validation.Validate(&reqBody)

	result, err := handler.UserUsecase.Update(id, &reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)

}

func (handler UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)

	err = handler.UserUsecase.Delete(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}
