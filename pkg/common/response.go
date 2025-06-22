package common

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"` // for pagination
}

func Success(c *fiber.Ctx, data interface{}, message string, meta ...interface{}) error {
	resp := Response{
		Code:    fiber.StatusOK,
		Message: message,
		Data:    data,
	}
	if len(meta) > 0 {
		resp.Meta = meta[0]
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Code:    fiber.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(Response{
		Code:    status,
		Message: err.Error(),
	})
}

func BadRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Code:    fiber.StatusBadRequest,
		Message: err.Error(),
	})
}

func NotFound(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Code:    fiber.StatusNotFound,
		Message: err.Error(),
	})
}
