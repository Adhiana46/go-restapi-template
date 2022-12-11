package main

import (
	"time"

	"github.com/Adhiana46/go-restapi-template/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func routes() *fiber.App {
	r := fiber.New(fiber.Config{
		ErrorHandler: handleError,
	})

	// Logger
	r.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		Done: func(c *fiber.Ctx, logString []byte) {
			if c.Response().StatusCode() != fiber.StatusOK {
				// reporter.SendToSlack(logString)
			}
		},
	}))

	// Handle Panic
	r.Use(func(c *fiber.Ctx) error {
		defer handlePanic(c)
		return c.Next()
	})

	api := r.Group("/api/v1")

	// Register Handlers
	handlers.
		NewActivityGroupHandler(svcActivityGroup).
		RegisterRoutes(api.Group("/activity-group"))
	handlers.
		NewTodoItemHandler(svcTodoItem).
		RegisterRoutes(api.Group("/activity-group/:activity_uuid/todo-items"))

	return r
}
