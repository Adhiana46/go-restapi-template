package repository

import (
	"fmt"

	"github.com/Adhiana46/go-restapi-template/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type TodoItemRepository interface {
	BeginTx() *sqlx.Tx

	FindByUuid(uuid string) (*entity.TodoItem, error)
	FindByUuidTx(tx *sqlx.Tx, uuid string) (*entity.TodoItem, error)
	FetchAll(page int, limit int, sorts map[string]string, activityId int, filter string) ([]*entity.TodoItem, error)
	CountAll(activityId int, filter string) (int, error)
	Store(tx *sqlx.Tx, e *entity.TodoItem) (*entity.TodoItem, error)
	Update(tx *sqlx.Tx, e *entity.TodoItem) (*entity.TodoItem, error)
	Delete(tx *sqlx.Tx, e *entity.TodoItem) error
}

type todoItemRepositoryPostgres struct {
	db *sqlx.DB
}

func (a *todoItemRepositoryPostgres) TableName() string {
	return "todo_item"
}

func (a *todoItemRepositoryPostgres) PrimaryField() string {
	return "id"
}

func NewPostgresTodoItemRepository(db *sqlx.DB) TodoItemRepository {
	return &todoItemRepositoryPostgres{
		db: db,
	}
}

func (r *todoItemRepositoryPostgres) BeginTx() *sqlx.Tx {
	return r.db.MustBegin()
}

func (r *todoItemRepositoryPostgres) FindByUuid(uuid string) (*entity.TodoItem, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("*").
		From(r.TableName()).
		Where(sq.Eq{"uuid": uuid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := entity.TodoItem{}
	err = r.db.Get(&row, sql, args...)
	if err != nil {
		return nil, err
	}

	if row.Activity != nil && row.Activity.ID == 0 {
		row.Activity = nil
	}

	return &row, nil
}

func (r *todoItemRepositoryPostgres) FindByUuidTx(tx *sqlx.Tx, uuid string) (*entity.TodoItem, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("*").
		From(r.TableName()).
		Where(sq.Eq{"uuid": uuid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := entity.TodoItem{}
	err = tx.Get(&row, sql, args...)
	if err != nil {
		return nil, err
	}

	if row.Activity != nil && row.Activity.ID == 0 {
		row.Activity = nil
	}

	return &row, nil
}

func (r *todoItemRepositoryPostgres) FetchAll(page int, limit int, sorts map[string]string, activityId int, filter string) ([]*entity.TodoItem, error) {
	offset := (page - 1) * limit

	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("*").
		From(r.TableName()).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if activityId != 0 {
		queryBuilder = queryBuilder.Where("activity_id = ?", activityId)
	}

	if filter != "" {
		queryBuilder = queryBuilder.Where("LOWER(name) LIKE ?", fmt.Sprint("%", filter, "%"))
	}

	if len(sorts) > 0 {
		for sortField, sortDir := range sorts {
			queryBuilder = queryBuilder.OrderBy(sortField + " " + sortDir)
		}
	}

	sql, args, err := queryBuilder.ToSql()

	if err != nil {
		return nil, err
	}

	rows := []*entity.TodoItem{}
	err = r.db.Select(&rows, sql, args...)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row.Activity != nil && row.Activity.ID == 0 {
			row.Activity = nil
		}
	}

	return rows, nil
}

func (r *todoItemRepositoryPostgres) CountAll(activityId int, filter string) (int, error) {
	total := 0

	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("COUNT(id) AS total").
		From(r.TableName())

	if activityId != 0 {
		queryBuilder = queryBuilder.Where("activity_id = ?", activityId)
	}

	if filter != "" {
		queryBuilder = queryBuilder.Where("LOWER(name) LIKE ?", fmt.Sprint("%", filter, "%"))
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := r.db.Queryx(sql, args...)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}

	return total, nil
}

func (r *todoItemRepositoryPostgres) Store(tx *sqlx.Tx, e *entity.TodoItem) (*entity.TodoItem, error) {
	values := map[string]interface{}{
		"uuid":        e.Uuid,
		"activity_id": e.ActivityID,
		"name":        e.Name,
		"description": e.Description,
		"created_at":  e.CreatedAt,
		"updated_at":  e.UpdatedAt,
	}

	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Insert(r.TableName()).
		SetMap(values).
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	return r.FindByUuidTx(tx, e.Uuid)
}

func (r *todoItemRepositoryPostgres) Update(tx *sqlx.Tx, e *entity.TodoItem) (*entity.TodoItem, error) {
	values := map[string]interface{}{
		"activity_id": e.ActivityID,
		"name":        e.Name,
		"description": e.Description,
		"updated_at":  e.UpdatedAt,
	}

	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update(r.TableName()).
		SetMap(values).
		Where(sq.Eq{"id": e.ID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	if e.Activity != nil && e.Activity.ID == 0 {
		e.Activity = nil
	}

	return e, nil
}

func (r *todoItemRepositoryPostgres) Delete(tx *sqlx.Tx, e *entity.TodoItem) error {
	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Delete(r.TableName()).
		Where(sq.Eq{"id": e.ID}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
