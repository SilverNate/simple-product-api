package middleware

import (
	"github.com/gofiber/fiber/v2"
	"simple-product-api/pkg/common"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := err.Error()

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(common.Response{
		Code:    code,
		Message: message,
	})
}
