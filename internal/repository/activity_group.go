package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Adhiana46/go-restapi-template/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type ActivityGroupRepository interface {
	BeginTx(ctx context.Context) *sqlx.Tx

	FindById(ctx context.Context, id int) (*entity.ActivityGroup, error)
	FindByUuid(ctx context.Context, uuid string) (*entity.ActivityGroup, error)
	FindByUuidTx(ctx context.Context, tx *sqlx.Tx, uuid string) (*entity.ActivityGroup, error)
	FetchAll(ctx context.Context, page int, limit int, sorts map[string]string, filter string) ([]*entity.ActivityGroup, error)
	CountAll(ctx context.Context, filter string) (int, error)
	Store(ctx context.Context, tx *sqlx.Tx, e *entity.ActivityGroup) (*entity.ActivityGroup, error)
	Update(ctx context.Context, tx *sqlx.Tx, e *entity.ActivityGroup) (*entity.ActivityGroup, error)
	Delete(ctx context.Context, tx *sqlx.Tx, e *entity.ActivityGroup) error
}

type activityGroupRepositoryPostgres struct {
	db *sqlx.DB
}

func (r *activityGroupRepositoryPostgres) TableName() string {
	return "activity_group"
}

func (r *activityGroupRepositoryPostgres) PrimaryField() string {
	return "id"
}

func NewPostgresActivityGroupRepository(db *sqlx.DB) ActivityGroupRepository {
	return &activityGroupRepositoryPostgres{
		db: db,
	}
}

func (r *activityGroupRepositoryPostgres) BeginTx(ctx context.Context) *sqlx.Tx {
	return r.db.MustBeginTx(ctx, &sql.TxOptions{})
}

func (r *activityGroupRepositoryPostgres) FindById(ctx context.Context, id int) (*entity.ActivityGroup, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("*").
		From(r.TableName()).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := entity.ActivityGroup{}
	err = r.db.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *activityGroupRepositoryPostgres) FindByUuid(ctx context.Context, uuid string) (*entity.ActivityGroup, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("*").
		From(r.TableName()).
		Where(sq.Eq{"uuid": uuid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := entity.ActivityGroup{}
	err = r.db.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *activityGroupRepositoryPostgres) FindByUuidTx(ctx context.Context, tx *sqlx.Tx, uuid string) (*entity.ActivityGroup, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("*").
		From(r.TableName()).
		Where(sq.Eq{"uuid": uuid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := entity.ActivityGroup{}
	err = tx.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *activityGroupRepositoryPostgres) FetchAll(ctx context.Context, page int, limit int, sorts map[string]string, filter string) ([]*entity.ActivityGroup, error) {
	offset := (page - 1) * limit

	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("*").
		From(r.TableName()).
		Limit(uint64(limit)).
		Offset(uint64(offset))

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

	rows := []*entity.ActivityGroup{}
	err = r.db.SelectContext(ctx, &rows, sql, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *activityGroupRepositoryPostgres) CountAll(ctx context.Context, filter string) (int, error) {
	total := 0

	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("COUNT(id) AS total").
		From(r.TableName())

	if filter != "" {
		queryBuilder = queryBuilder.Where("LOWER(name) LIKE ?", fmt.Sprint("%", filter, "%"))
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := r.db.QueryxContext(ctx, sql, args...)
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

func (r *activityGroupRepositoryPostgres) Store(ctx context.Context, tx *sqlx.Tx, e *entity.ActivityGroup) (*entity.ActivityGroup, error) {
	values := map[string]interface{}{
		"uuid":        e.Uuid,
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

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return r.FindByUuidTx(ctx, tx, e.Uuid)
}

func (r *activityGroupRepositoryPostgres) Update(ctx context.Context, tx *sqlx.Tx, e *entity.ActivityGroup) (*entity.ActivityGroup, error) {
	values := map[string]interface{}{
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

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *activityGroupRepositoryPostgres) Delete(ctx context.Context, tx *sqlx.Tx, e *entity.ActivityGroup) error {
	// Build SQL
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Delete(r.TableName()).
		Where(sq.Eq{"id": e.ID}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
