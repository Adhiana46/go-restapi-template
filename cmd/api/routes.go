package main

import (
	"log"

	"github.com/Adhiana46/go-restapi-template/handlers"
	"github.com/gofiber/fiber/v2"
)

func routes() *fiber.App {
	r := fiber.New(fiber.Config{
		ErrorHandler: handleError,
	})

	// Handle Panic
	r.Use(func(c *fiber.Ctx) error {
		log.Println("Handle Panic Middleware")

		defer handlePanic(c)

		return c.Next()
	})

	api := r.Group("/api/v1")

	// Register Handlers
	handlers.
		NewActivityGroupHandler(svcActivityGroup).
		RegisterRoutes(api.Group("/activity_group"))

	return r
}
