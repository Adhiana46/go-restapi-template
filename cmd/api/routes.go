package main

import (
	"time"

	httpTransport "github.com/Adhiana46/go-restapi-template/transport/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func httpRoutes() *fiber.App {
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
	httpTransport.
		NewActivityGroupHandler(svcActivityGroup).
		RegisterRoutes(api.Group("/activity-group"))
	httpTransport.
		NewTodoItemHandler(svcTodoItem).
		RegisterRoutes(api.Group("/activity-group/:activity_uuid/todo-items"))

	return r
}
