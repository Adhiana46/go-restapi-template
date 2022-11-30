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
				if err := shouldBind(c, &req); err != nil {
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
				if err := shouldBind(c, &req); err != nil {
					panic(err)
				}

				activityGroupList, pagination, err := svcActivityGroup.FetchAll(req)
				if err != nil {
					return handleError(c, err)
				}

				resp := dto.ActivityGroupToResponseList(activityGroupList)

				return c.JSON(responsePkg.JsonSuccess(http.StatusOK, "", resp, pagination))
			})
			api.Post("activity_group", func(c *fiber.Ctx) error {
				defer handlePanic(c)

				req := dto.ActivityGroupCreateRequest{}
				if err := shouldBind(c, &req); err != nil {
					panic(err)
				}

				activityGroup, err := svcActivityGroup.Create(req)
				if err != nil {
					return handleError(c, err)
				}

				resp := dto.ActivityGroupToResponse(activityGroup)

				return c.JSON(responsePkg.JsonSuccess(http.StatusOK, "", resp, nil))
			})
			api.Put("activity_group/:uuid", func(c *fiber.Ctx) error {
				defer handlePanic(c)

				req := dto.ActivityGroupUpdateRequest{}
				if err := shouldBind(c, &req); err != nil {
					panic(err)
				}

				activityGroup, err := svcActivityGroup.Update(req)
				if err != nil {
					return handleError(c, err)
				}

				resp := dto.ActivityGroupToResponse(activityGroup)

				return c.JSON(responsePkg.JsonSuccess(http.StatusOK, "", resp, nil))
			})
			api.Delete("activity_group/:uuid", func(c *fiber.Ctx) error {
				defer handlePanic(c)

				req := dto.ActivityGroupUuidRequest{}
				if err := shouldBind(c, &req); err != nil {
					panic(err)
				}

				err := svcActivityGroup.Delete(req)
				if err != nil {
					return handleError(c, err)
				}

				return c.JSON(responsePkg.JsonSuccess(http.StatusNoContent, "", nil, nil))
			})
		}
	}

	return r
}
