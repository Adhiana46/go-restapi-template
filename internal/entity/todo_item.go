package entity

import "time"

type TodoItem struct {
	ID          int            `db:"id" json:"id"`
	Uuid        string         `db:"uuid" json:"uuid"`
	ActivityID  int            `db:"activity_id" json:"activity_id"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description" json:"description"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
	Activity    *ActivityGroup `json:"activity,omitempty"`
}
