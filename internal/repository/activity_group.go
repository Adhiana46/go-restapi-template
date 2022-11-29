package repository

import (
	"strconv"
	"strings"

	"github.com/Adhiana46/go-restapi-template/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ActivityGroupRepository interface {
	FindByUuid(uuid string) (*entity.ActivityGroup, error)
	FetchAll(page int, limit int, sorts map[string]string, filter string) ([]*entity.ActivityGroup, error)
	CountAll(filter string) (int, error)
	// Store(e *entity.ActivityGroup) (*entity.ActivityGroup, error)
	// Update(id int, e *entity.ActivityGroup) (*entity.ActivityGroup, error)
	// DeleteById(id int) (bool, error)
	// Delete(e *entity.ActivityGroup) (bool, error)
}

type activityGroupRepositoryPostgres struct {
	db *sqlx.DB
}

func NewPostgresActivityGroupRepository(db *sqlx.DB) ActivityGroupRepository {
	return &activityGroupRepositoryPostgres{
		db: db,
	}
}

func (r *activityGroupRepositoryPostgres) FindByUuid(uuid string) (*entity.ActivityGroup, error) {
	row := entity.ActivityGroup{}
	err := r.db.Get(&row, "SELECT * FROM activity_group WHERE uuid = $1", uuid)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *activityGroupRepositoryPostgres) FetchAll(page int, limit int, sorts map[string]string, filter string) ([]*entity.ActivityGroup, error) {
	offset := (page - 1) * limit

	// Build SQL
	sqlWhere := ""
	sqlOrders := ""
	sqlLimit := "LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)

	if filter != "" {
		sqlWhere = "WHERE LOWER(name) LIKE :filter"
	}

	if len(sorts) > 0 {
		aOrders := []string{}
		for sortField, sortDir := range sorts {
			aOrders = append(aOrders, sortField+" "+sortDir)
		}

		sqlOrders = "ORDER BY " + strings.Join(aOrders, ",")
	}

	sql := "SELECT * FROM activity_group " + sqlWhere + " " + sqlOrders + " " + sqlLimit

	rows := []*entity.ActivityGroup{}
	result, err := r.db.NamedQuery(sql, map[string]interface{}{"filter": "%" + filter + "%"})
	if err != nil {
		return nil, err
	}
	for result.Next() {
		row := entity.ActivityGroup{}
		err = result.StructScan(&row)
		if err != nil {
			return nil, err
		}

		rows = append(rows, &row)
	}

	return rows, nil
}

func (r *activityGroupRepositoryPostgres) CountAll(filter string) (int, error) {
	total := 0

	sqlWhere := ""

	if filter != "" {
		sqlWhere = "WHERE LOWER(name) LIKE :filter"
	}

	sql := "SELECT COUNT(id) AS total FROM activity_group " + sqlWhere

	var err error = nil
	rows, err := r.db.NamedQuery(sql, map[string]interface{}{"filter": "%" + filter + "%"})
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
