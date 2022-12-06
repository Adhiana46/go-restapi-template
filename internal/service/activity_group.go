package service

import (
	"context"
	"math"
	"time"

	"github.com/Adhiana46/go-restapi-template/internal/dto"
	"github.com/Adhiana46/go-restapi-template/internal/entity"
	"github.com/Adhiana46/go-restapi-template/internal/repository"
	parserPkg "github.com/Adhiana46/go-restapi-template/pkg/parser"
	responsePkg "github.com/Adhiana46/go-restapi-template/pkg/response"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ActivityGroupService interface {
	FindByUuid(ctx context.Context, req dto.ActivityGroupUuidRequest) (*entity.ActivityGroup, error)
	FetchAll(ctx context.Context, req dto.ActivityGroupFetchRequest) ([]*entity.ActivityGroup, *responsePkg.Pagination, error)
	Create(ctx context.Context, req dto.ActivityGroupCreateRequest) (*entity.ActivityGroup, error)
	Update(ctx context.Context, req dto.ActivityGroupUpdateRequest) (*entity.ActivityGroup, error)
	Delete(ctx context.Context, req dto.ActivityGroupUuidRequest) error
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

func (s *activityGroupService) FindByUuid(ctx context.Context, req dto.ActivityGroupUuidRequest) (*entity.ActivityGroup, error) {
	span := sentry.StartSpan(ctx, "activityGroupService.FindByUuid")
	defer span.Finish()

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	activityGroup, err := s.repo.FindByUuid(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	return activityGroup, nil
}

func (s *activityGroupService) FetchAll(ctx context.Context, req dto.ActivityGroupFetchRequest) ([]*entity.ActivityGroup, *responsePkg.Pagination, error) {
	span := sentry.StartSpan(ctx, "activityGroupService.FetchAll")
	defer span.Finish()

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

	totalRows, err := s.repo.CountAll(ctx, req.Filter)
	if err != nil {
		return nil, nil, err
	}

	sorts, err := parserPkg.QuerySortToMap(req.SortBy)
	if err != nil {
		return nil, nil, err
	}

	activityGroupList, err := s.repo.FetchAll(ctx, req.Page, req.Limit, sorts, req.Filter)
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

func (s *activityGroupService) Create(ctx context.Context, req dto.ActivityGroupCreateRequest) (*entity.ActivityGroup, error) {
	span := sentry.StartSpan(ctx, "activityGroupService.Create")
	defer span.Finish()

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	ent := &entity.ActivityGroup{
		Uuid:        uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// begin transaction
	tx := s.repo.BeginTx()
	insertedRow, err := s.repo.Store(ctx, tx, ent)

	// if error rollback, commit otherwise
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
	}

	return insertedRow, nil
}

func (s *activityGroupService) Update(ctx context.Context, req dto.ActivityGroupUpdateRequest) (*entity.ActivityGroup, error) {
	span := sentry.StartSpan(ctx, "activityGroupService.Update")
	defer span.Finish()

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	ent, err := s.repo.FindByUuid(ctx, req.Uuid)
	if err != nil {
		return ent, err
	}

	// Update values
	ent.Name = req.Name
	ent.Description = req.Description
	ent.UpdatedAt = time.Now()

	// begin transaction
	tx := s.repo.BeginTx()
	updatedRow, err := s.repo.Update(ctx, tx, ent)

	// if error rollback, commit otherwise
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
	}

	return updatedRow, nil
}

func (s *activityGroupService) Delete(ctx context.Context, req dto.ActivityGroupUuidRequest) error {
	span := sentry.StartSpan(ctx, "activityGroupService.Delete")
	defer span.Finish()

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	ent, err := s.repo.FindByUuid(ctx, req.Uuid)
	if err != nil {
		return err
	}

	// begin transaction
	tx := s.repo.BeginTx()
	err = s.repo.Delete(ctx, tx, ent)

	// if error rollback, commit otherwise
	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
	}

	return nil
}
