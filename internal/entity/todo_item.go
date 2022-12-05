package entity

import "time"

type TodoItem struct {
	ID          int       `db:"id"`
	Uuid        string    `db:"uuid"`
	ActivityID  int       `db:"activity_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Activity    *ActivityGroup
}
