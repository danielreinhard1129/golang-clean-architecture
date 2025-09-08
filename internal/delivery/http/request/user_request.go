package request

import (
	"strconv"

	"github.com/danielreinhard1129/fiber-clean-arch/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type UserFindAllRequest struct {
	Page    int    `json:"page" validate:"gte=1"`
	Limit   int    `json:"limit" validate:"gte=1,lte=100"`
	Search  string `json:"search" validate:"omitempty,min=1"`
	OrderBy string `json:"orderBy" validate:"omitempty,min=1"`
	Sort    string `json:"sort" validate:"omitempty,min=1"`
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

func ParseAndValidate(c *fiber.Ctx) UserFindAllRequest {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	orderBy := c.Query("orderBy", "created_at")
	sort := c.Query("sort", "desc")

	req := UserFindAllRequest{
		Page:    page,
		Limit:   limit,
		Search:  search,
		OrderBy: orderBy,
		Sort:    sort,
	}

	validation.Validate(&req)

	return req
}
