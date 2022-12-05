package service

import (
	"math"
	"time"

	"github.com/Adhiana46/go-restapi-template/internal/dto"
	"github.com/Adhiana46/go-restapi-template/internal/entity"
	"github.com/Adhiana46/go-restapi-template/internal/repository"
	parserPkg "github.com/Adhiana46/go-restapi-template/pkg/parser"
	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TodoItemService interface {
	FindByUuid(req dto.TodoItemUuidRequest) (*entity.TodoItem, error)
	FetchAll(req dto.TodoItemFetchRequest) ([]*entity.TodoItem, *responsePkg.Pagination, error)
	Create(req dto.TodoItemCreateRequest) (*entity.TodoItem, error)
	Update(req dto.TodoItemUpdateRequest) (*entity.TodoItem, error)
	Delete(req dto.TodoItemUuidRequest) error
}

type todoItemService struct {
	validate     *validator.Validate
	repo         repository.TodoItemRepository
	repoActivity repository.ActivityGroupRepository
}

func NewTodoItemService(validate *validator.Validate, repo repository.TodoItemRepository, repoActivity repository.ActivityGroupRepository) TodoItemService {
	return &todoItemService{
		validate:     validate,
		repo:         repo,
		repoActivity: repoActivity,
	}
}

func (s *todoItemService) FindByUuid(req dto.TodoItemUuidRequest) (*entity.TodoItem, error) {
	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	todoItem, err := s.repo.FindByUuid(req.Uuid)
	if err != nil {
		return nil, err
	}

	activity, err := s.repoActivity.FindById(todoItem.ActivityID)
	if err != nil {
		return nil, err
	}
	todoItem.Activity = activity

	return todoItem, nil
}

func (s *todoItemService) FetchAll(req dto.TodoItemFetchRequest) ([]*entity.TodoItem, *responsePkg.Pagination, error) {
	var err error

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

	activity := &entity.ActivityGroup{}
	if req.ActivityUuid != "" {
		activity, err = s.repoActivity.FindByUuid(req.ActivityUuid)
		if err != nil {
			return nil, nil, err
		}
	}

	totalRows, err := s.repo.CountAll(activity.ID, req.Filter)
	if err != nil {
		return nil, nil, err
	}

	sorts, err := parserPkg.QuerySortToMap(req.SortBy)
	if err != nil {
		return nil, nil, err
	}

	todoItemList, err := s.repo.FetchAll(req.Page, req.Limit, sorts, activity.ID, req.Filter)
	if err != nil {
		return nil, nil, err
	}

	// Create Pagination
	pagination := responsePkg.Pagination{
		CurrentPage: req.Page,
		Total:       totalRows,
		Size:        len(todoItemList),
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(req.Limit))),
	}

	return todoItemList, &pagination, nil
}

func (s *todoItemService) Create(req dto.TodoItemCreateRequest) (*entity.TodoItem, error) {
	var err error

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	activity := &entity.ActivityGroup{}
	if req.ActivityUuid != "" {
		activity, err = s.repoActivity.FindByUuid(req.ActivityUuid)
		if err != nil {
			return nil, err
		}
	}

	ent := &entity.TodoItem{
		Uuid:        uuid.NewString(),
		ActivityID:  activity.ID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// begin transaction
	tx := s.repo.BeginTx()
	insertedRow, err := s.repo.Store(tx, ent)

	// if error rollback, commit otherwise
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
	}

	return insertedRow, nil
}

func (s *todoItemService) Update(req dto.TodoItemUpdateRequest) (*entity.TodoItem, error) {
	var err error

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	ent, err := s.repo.FindByUuid(req.Uuid)
	if err != nil {
		return ent, err
	}

	activity := &entity.ActivityGroup{}
	if req.ActivityUuid != "" {
		activity, err = s.repoActivity.FindByUuid(req.ActivityUuid)
		if err != nil {
			return nil, err
		}
	}

	// Update values
	ent.ActivityID = activity.ID
	ent.Name = req.Name
	ent.Description = req.Description
	ent.UpdatedAt = time.Now()

	// begin transaction
	tx := s.repo.BeginTx()
	updatedRow, err := s.repo.Update(tx, ent)

	// if error rollback, commit otherwise
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
	}

	return updatedRow, nil
}

func (s *todoItemService) Delete(req dto.TodoItemUuidRequest) error {
	// Validate
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	ent, err := s.repo.FindByUuid(req.Uuid)
	if err != nil {
		return err
	}

	// begin transaction
	tx := s.repo.BeginTx()
	err = s.repo.Delete(tx, ent)

	// if error rollback, commit otherwise
	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
	}

	return nil
}
