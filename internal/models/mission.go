package models

import (
	"time"

	"gorm.io/gorm"
)

type Mission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CatID       *uint          `json:"cat_id" gorm:"index"`
	Cat         *SpyCat        `json:"cat,omitempty" gorm:"foreignKey:CatID"`
	IsCompleted bool           `json:"is_completed" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Targets     []Target       `json:"targets,omitempty" gorm:"foreignKey:MissionID"`
}

type Target struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MissionID   uint           `json:"mission_id" gorm:"not null;index"`
	Mission     *Mission       `json:"mission,omitempty" gorm:"foreignKey:MissionID"`
	Name        string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Country     string         `json:"country" gorm:"not null" validate:"required,min=2,max=100"`
	Notes       string         `json:"notes" gorm:"type:text"`
	IsCompleted bool           `json:"is_completed" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateMissionRequest struct {
	CatID   *uint                 `json:"cat_id"`
	Targets []CreateTargetRequest `json:"targets" validate:"required,min=1,max=3,dive"`
}

type CreateTargetRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Country string `json:"country" validate:"required,min=2,max=100"`
}

type UpdateMissionRequest struct {
	CatID *uint `json:"cat_id"`
}

type AddTargetRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Country string `json:"country" validate:"required,min=2,max=100"`
}

type UpdateTargetRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Country string `json:"country" validate:"required,min=2,max=100"`
}

type UpdateTargetNotesRequest struct {
	Notes string `json:"notes" validate:"required"`
}
