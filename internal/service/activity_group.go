package service

import (
	"math"

	"github.com/Adhiana46/go-restapi-template/internal/dto"
	"github.com/Adhiana46/go-restapi-template/internal/entity"
	"github.com/Adhiana46/go-restapi-template/internal/repository"
	parserPkg "github.com/Adhiana46/go-restapi-template/pkg/parser"
	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	"github.com/go-playground/validator/v10"
)

type ActivityGroupService interface {
	FindByUuid(req dto.ActivityGroupUuidRequest) (*entity.ActivityGroup, error)
	FetchAll(req dto.ActivityGroupFetchRequest) ([]*entity.ActivityGroup, *responsePkg.Pagination, error)
	// Create(req dto.CreateActivityGroup) (*entity.ActivityGroup, error)
	// Update(req dto.UpdateActivityGroup) (*entity.ActivityGroup, error)
	// DeleteById(req dto.ActivityGroupId) error
}

type activityGroupService struct {
	validate *validator.Validate
	repo     repository.ActivityGroupRepository
}

func NewActivityGroupService(validate *validator.Validate, repo repository.ActivityGroupRepository) ActivityGroupService {
	return &activityGroupService{
		validate: validate,
		repo:     repo,
	}
}

func (s *activityGroupService) FindByUuid(req dto.ActivityGroupUuidRequest) (*entity.ActivityGroup, error) {
	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	activityGroup, err := s.repo.FindByUuid(req.Uuid)
	if err != nil {
		return nil, err
	}

	return activityGroup, nil
}

func (s *activityGroupService) FetchAll(req dto.ActivityGroupFetchRequest) ([]*entity.ActivityGroup, *responsePkg.Pagination, error) {
	// Set Default Value
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.SortBy == "" {
		req.SortBy = "name.asc"
	}

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, nil, err
	}

	totalRows, err := s.repo.CountAll(req.Filter)
	if err != nil {
		return nil, nil, err
	}

	sorts, err := parserPkg.QuerySortToMap(req.SortBy)
	if err != nil {
		return nil, nil, err
	}

	activityGroupList, err := s.repo.FetchAll(req.Page, req.Limit, sorts, req.Filter)
	if err != nil {
		return nil, nil, err
	}

	// Create Pagination
	pagination := responsePkg.Pagination{
		CurrentPage: req.Page,
		Total:       totalRows,
		Size:        len(activityGroupList),
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(req.Limit))),
	}

	return activityGroupList, &pagination, nil
}
