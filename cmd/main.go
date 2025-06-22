package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "simple-product-api/docs"
	"simple-product-api/pkg/common"
	"simple-product-api/pkg/config"
	"simple-product-api/pkg/di"
	middleware "simple-product-api/pkg/midlleware"
	"time"
)

// @title Simple Product API
// @version 1.0
// @description REST API for managing products.
// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	cfg := config.Load()
	logrus.Info("Server starting...")

	handler, err := di.InitializeHandler(cfg)
	if err != nil {
		logrus.Fatalf("failed to initialize product handler: %v", err)
	}
	logrus.Info("intialize product handler successfully")

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	logrus.Info("retry logic is set")
	app.Use(middleware.RetryWithTimeout(3*time.Second, 2))
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	logrus.Info("rate limiting is set")
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(common.Response{
				Code:    fiber.StatusTooManyRequests,
				Message: "Too many requests. Please try again later.",
			})
		},
	}))

	api := app.Group("/api/v1")
	handler.Register(api.Group("/products"))

	logrus.Fatal(app.Listen(":8080"))
}
