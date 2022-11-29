package entity

import "time"

type ActivityGroup struct {
	ID          int       `db:"id"`
	Uuid        string    `db:"uuid"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
