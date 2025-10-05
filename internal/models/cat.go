package models

import (
	"time"

	"gorm.io/gorm"
)

type SpyCat struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	YearsExperience int            `json:"years_experience" gorm:"not null" validate:"required,min=0,max=50"`
	Breed           string         `json:"breed" gorm:"not null" validate:"required"`
	Salary          float64        `json:"salary" gorm:"not null" validate:"required,min=0"`
	IsAvailable     bool           `json:"is_available" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateCatRequest struct {
	Name            string  `json:"name" validate:"required,min=2,max=100"`
	YearsExperience int     `json:"years_experience" validate:"required,min=0,max=50"`
	Breed           string  `json:"breed" validate:"required"`
	Salary          float64 `json:"salary" validate:"required,min=0"`
}

type UpdateCatRequest struct {
	Salary float64 `json:"salary" validate:"required,min=0"`
}
