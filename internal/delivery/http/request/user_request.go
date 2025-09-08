package request

import (
	"strconv"

	"github.com/danielreinhard1129/fiber-clean-arch/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type UserFindAllRequest struct {
	Page   int    `json:"page" validate:"gte=1"`
	Limit  int    `json:"limit" validate:"gte=1,lte=100"`
	Search string `json:"search" validate:"omitempty,min=1"`
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func ParseAndValidate(c *fiber.Ctx) UserFindAllRequest {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	req := UserFindAllRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	validation.Validate(&req)

	return req
}
