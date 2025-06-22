package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RetryWithTimeout(timeout time.Duration, maxRetries int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		for i := 0; i <= maxRetries; i++ {
			c.SetUserContext(ctx)

			err = c.Next()

			if err == nil || ctx.Err() == context.DeadlineExceeded {
				break
			}
		}

		if err != nil {
			return fiber.NewError(fiber.StatusRequestTimeout, "Request failed or timed out")
		}
		return nil
	}
}
