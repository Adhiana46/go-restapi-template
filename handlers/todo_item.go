package handlers

import (
	"net/http"

	"github.com/Adhiana46/go-restapi-template/internal/dto"
	"github.com/Adhiana46/go-restapi-template/internal/service"
	parserPkg "github.com/Adhiana46/go-restapi-template/pkg/parser"
	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type TodoItemHandler interface {
	RegisterRoutes(r fiber.Router) TodoItemHandler

	findByUuid() func(c *fiber.Ctx) error
	fetchAll() func(c *fiber.Ctx) error
	create() func(c *fiber.Ctx) error
	update() func(c *fiber.Ctx) error
	delete() func(c *fiber.Ctx) error
}

type todoItemHandler struct {
	svcTodoItem service.TodoItemService
}

func NewTodoItemHandler(svcTodoItem service.TodoItemService) TodoItemHandler {
	return &todoItemHandler{
		svcTodoItem: svcTodoItem,
	}
}

func (h *todoItemHandler) RegisterRoutes(r fiber.Router) TodoItemHandler {
	r.Get("/:uuid", h.findByUuid())
	r.Get("/", h.fetchAll())
	r.Post("/", h.create())
	r.Put("/:uuid", h.update())
	r.Delete("/:uuid", h.delete())

	return h
}

func (h *todoItemHandler) findByUuid() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.TodoItemUuidRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		todoItem, err := h.svcTodoItem.FindByUuid(req)
		if err != nil {
			return err
		}

		resp := dto.TodoItemToResponse(todoItem)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, nil))
	}
}

func (h *todoItemHandler) fetchAll() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.TodoItemFetchRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		req.ActivityUuid = c.Params("activity_uuid")

		todoItemList, pagination, err := h.svcTodoItem.FetchAll(req)
		if err != nil {
			return err
		}

		resp := dto.TodoItemToResponseList(todoItemList)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, pagination))
	}
}

func (h *todoItemHandler) create() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.TodoItemCreateRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		req.ActivityUuid = c.Params("activity_uuid")

		todoItem, err := h.svcTodoItem.Create(req)
		if err != nil {
			return err
		}

		resp := dto.TodoItemToResponse(todoItem)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, nil))
	}
}

func (h *todoItemHandler) update() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.TodoItemUpdateRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		req.ActivityUuid = c.Params("activity_uuid")

		todoItem, err := h.svcTodoItem.Update(req)
		if err != nil {
			return err
		}

		resp := dto.TodoItemToResponse(todoItem)

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", resp, nil))
	}
}

func (h *todoItemHandler) delete() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := dto.TodoItemUuidRequest{}
		if err := parserPkg.FiberShouldBindRequest(c, &req); err != nil {
			panic(err)
		}

		err := h.svcTodoItem.Delete(req)
		if err != nil {
			return err
		}

		statusCode := http.StatusOK
		return c.Status(statusCode).JSON(responsePkg.JsonSuccess(statusCode, "", nil, nil))
	}
}
