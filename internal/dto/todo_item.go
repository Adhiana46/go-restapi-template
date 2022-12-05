package dto

import (
	"time"

	"github.com/Adhiana46/go-restapi-template/internal/entity"
)

func TodoItemToResponse(e *entity.TodoItem) *TodoItemResponse {
	resp := &TodoItemResponse{
		Uuid:        e.Uuid,
		ActivityID:  e.ActivityID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	if e.Activity != nil {
		resp.Activity = ActivityGroupToResponse(e.Activity)
	}

	return resp
}

func TodoItemToResponseList(ents []*entity.TodoItem) []*TodoItemResponse {
	respList := []*TodoItemResponse{}

	for _, e := range ents {
		respList = append(respList, TodoItemToResponse(e))
	}

	return respList
}

type TodoItemResponse struct {
	Uuid        string                 `json:"uuid"`
	ActivityID  int                    `json:"activity_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Activity    *ActivityGroupResponse `json:"activity,omitempty"`
}

type TodoItemUuidRequest struct {
	Uuid string `uri:"uuid" validate:"required"`
}

type TodoItemFetchRequest struct {
	Page         int    `query="page" validate:"numeric,min=1"`
	Limit        int    `query="limit" validate:"numeric,min=1,max=200"`
	SortBy       string `query="sortBy" validate:""`
	ActivityUuid string `uri:"activity_uuid" query:"activity_uuid"`
	Filter       string `query="filter" validate:""`
}

type TodoItemCreateRequest struct {
	ActivityUuid string `json:"activity_uuid" uri:"activity_uuid" validate:"required"`
	Name         string `json:"name" validate:"required,min=3,max=100"`
	Description  string `json:"description" validate:""`
}

type TodoItemUpdateRequest struct {
	Uuid         string `uri:"uuid" validate:"required"`
	ActivityUuid string `json:"activity_uuid" uri:"activity_uuid" validate:"required"`
	Name         string `json:"name" validate:"required,min=3,max=100"`
	Description  string `json:"description" validate:""`
}
