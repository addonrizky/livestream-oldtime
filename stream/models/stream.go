package models

import "time"

type Stream struct {
	ID          int64      `json:"id"`
	Name        string     `json:"username" validate:"required"`
	CreatedBy   int64      `json:"created_by"`
	CreatedDate time.Time  `json:"created_date"`
	UpdatedBy   *int64     `json:"updated_by"`
	UpdatedDate *time.Time `json:"updated_date"`
}
