package models

import "time"

type Employee struct {
	ID           int        `gorm:"primaryKey" json:"id"`
	DepartmentID int        `gorm:"not null" json:"department_id"`
	FullName     string     `gorm:"not null" json:"full_name,omitempty"`
	Position     string     `gorm:"not null" json:"position,omitempty"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`

	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
