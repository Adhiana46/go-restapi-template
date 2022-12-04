package dto

import (
	"time"

	"github.com/Adhiana46/go-restapi-template/internal/entity"
)

func ActivityGroupToResponse(e *entity.ActivityGroup) *ActivityGroupResponse {
	return &ActivityGroupResponse{
		Uuid:        e.Uuid,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func ActivityGroupToResponseList(ents []*entity.ActivityGroup) []*ActivityGroupResponse {
	respList := []*ActivityGroupResponse{}

	for _, e := range ents {
		respList = append(respList, ActivityGroupToResponse(e))
	}

	return respList
}

type ActivityGroupResponse struct {
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ActivityGroupUuidRequest struct {
	Uuid string `uri:"uuid" validate:"required"`
}

type ActivityGroupFetchRequest struct {
	Page   int    `query="page" validate:"numeric,min=1"`
	Limit  int    `query="limit" validate:"numeric,min=1,max=200"`
	SortBy string `query="sortBy" validate:""`
	Filter string `query="filter" validate:""`
}

type ActivityGroupCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:""`
}

type ActivityGroupUpdateRequest struct {
	Uuid        string `uri:"uuid" validate:"required"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:""`
}
