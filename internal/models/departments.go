package models

import "time"

type Department struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name,omitempty"`
	ParentID *int   `gorm:"index" json:"parent_id,omitempty"`

	Parent   *Department `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"parent,omitempty"`
	Employee []Employee  `gorm:"foreignKey:DepartmentID;constraint:OnDelete:CASCADE" json:"employee,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
