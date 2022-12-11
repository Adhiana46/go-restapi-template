package http

import (
	"net/http"

	"github.com/Adhiana46/go-restapi-template/internal/dto"
	"github.com/Adhiana46/go-restapi-template/internal/service"
	parserPkg "github.com/Adhiana46/go-restapi-template/pkg/parser"
	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ActivityGroupHandler interface {
	RegisterRoutes(r fiber.Router) ActivityGroupHandler

	findByUuid() func(c *fiber.Ctx) error
	fetchAll() func(c *fiber.Ctx) error
	create() func(c *fiber.Ctx) error
	update() func(c *fiber.Ctx) error
	delete() func(c *fiber.Ctx) error
}

type activityGroupHandler struct {
	svcActivityGroup service.ActivityGroupService
}

func NewActivityGroupHandler(svcActivityGroup service.ActivityGroupService) ActivityGroupHandler {
	return &activityGroupHandler{
		svcActivityGroup: svcActivityGroup,
	}
}

func (h *activityGroupHandler) RegisterRoutes(r fiber.Router) ActivityGroupHandler {
	r.Get("/:uuid", h.findByUuid())
	r.Get("/", h.fetchAll())
	r.Post("/", h.create())
	r.Put("/:uuid", h.update())
	r.Delete("/:uuid", h.delete())

	return h
}

func (h *activityGroupHandler) findByUuid() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.ActivityGroupUuidRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		activityGroup, err := h.svcActivityGroup.FindByUuid(req)
		if err != nil {
			return err
		}

		resp := dto.ActivityGroupToResponse(activityGroup)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, nil))
	}
}

func (h *activityGroupHandler) fetchAll() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.ActivityGroupFetchRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		activityGroupList, pagination, err := h.svcActivityGroup.FetchAll(req)
		if err != nil {
			return err
		}

		resp := dto.ActivityGroupToResponseList(activityGroupList)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, pagination))
	}
}

func (h *activityGroupHandler) create() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.ActivityGroupCreateRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		activityGroup, err := h.svcActivityGroup.Create(req)
		if err != nil {
			return err
		}

		resp := dto.ActivityGroupToResponse(activityGroup)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, nil))
	}
}

func (h *activityGroupHandler) update() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.ActivityGroupUpdateRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		activityGroup, err := h.svcActivityGroup.Update(req)
		if err != nil {
			return err
		}

		resp := dto.ActivityGroupToResponse(activityGroup)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, nil))
	}
}

func (h *activityGroupHandler) delete() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.ActivityGroupUuidRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		err := h.svcActivityGroup.Delete(req)
		if err != nil {
			return err
		}

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", nil, nil))
	}
}
