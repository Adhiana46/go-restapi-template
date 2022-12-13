package entity

import "time"

type ActivityGroup struct {
	ID          int       `db:"id" json:"id"`
	Uuid        string    `db:"uuid" json:"uuid"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
