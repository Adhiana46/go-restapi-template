package main

import (
	"net/http"

	"github.com/Adhiana46/go-restapi-template/internal/dto"
	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func routes() *fiber.App {
	r := fiber.New()

	api := r.Group("/api/v1")
	{
		activityGroupRoutes := api.Group("/activity_group")
		{
			activityGroupRoutes.Get(":uuid", func(c *fiber.Ctx) error {
				defer handlePanic(c)

				req := dto.ActivityGroupUuidRequest{}

				if err := c.ParamsParser(&req); err != nil {
					panic(err)
				}

				activityGroup, err := svcActivityGroup.FindByUuid(req)
				if err != nil {
					return handleError(c, err)
				}

				resp := dto.ActivityGroupToResponse(activityGroup)

				return c.JSON(responsePkg.JsonSuccess(http.StatusOK, "", resp, nil))
			})
			api.Get("activity_group", func(c *fiber.Ctx) error {
				defer handlePanic(c)

				req := dto.ActivityGroupFetchRequest{}

				if err := c.QueryParser(&req); err != nil {
					panic(err)
				}

				activityGroupList, pagination, err := svcActivityGroup.FetchAll(req)
				if err != nil {
					return handleError(c, err)
				}

				resp := dto.ActivityGroupToResponseList(activityGroupList)

				return c.JSON(responsePkg.JsonSuccess(http.StatusOK, "", resp, pagination))
			})
			// api.Post("activity_group", activityGroupHandler.Create())
			// api.Put("activity_group/:id", activityGroupHandler.Update())
			// api.Delete("activity_group/:id", activityGroupHandler.Delete())
		}
	}

	return r
}
