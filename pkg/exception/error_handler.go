package exception

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	_, validationError := err.(ValidationError)
	if validationError {
		data := err.Error()
		var messages []map[string]interface{}
		errJson := json.Unmarshal([]byte(data), &messages)
		PanicLogging(errJson)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request",
			"data":    messages,
		})
	}

	_, notFoundError := err.(NotFoundError)
	if notFoundError {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"message": "Not Found",
			"data":    err.Error(),
		})
	}

	_, unauthorizedError := err.(UnauthorizedError)
	if unauthorizedError {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
			"data":    err.Error(),
		})
	}

	_, conflictError := err.(ConflictError)
	if conflictError {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"code":    409,
			"message": "Conflict",
			"data":    err.Error(),
		})
	}

	_, badRequestError := err.(BadRequestError)
	if badRequestError {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request",
			"data":    err.Error(),
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":    500,
		"message": "Error",
		"data":    err.Error(),
	})
}
