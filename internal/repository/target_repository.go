package repository

import (
	"spy-cat-agency/internal/models"

	"gorm.io/gorm"
)

type TargetRepository interface {
	Create(target *models.Target) error
	GetByID(id uint) (*models.Target, error)
	GetByMissionID(missionID uint) ([]models.Target, error)
	Update(target *models.Target) error
	Delete(id uint) error
	CompleteTarget(id uint) error
	UpdateNotes(id uint, notes string) error
	CountByMissionID(missionID uint) (int64, error)
}

type targetRepository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) TargetRepository {
	return &targetRepository{db: db}
}

func (r *targetRepository) Create(target *models.Target) error {
	return r.db.Create(target).Error
}

func (r *targetRepository) GetByID(id uint) (*models.Target, error) {
	var target models.Target
	err := r.db.Preload("Mission").First(&target, id).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *targetRepository) GetByMissionID(missionID uint) ([]models.Target, error) {
	var targets []models.Target
	err := r.db.Where("mission_id = ?", missionID).Find(&targets).Error
	return targets, err
}

func (r *targetRepository) Update(target *models.Target) error {
	return r.db.Save(target).Error
}

func (r *targetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Target{}, id).Error
}

func (r *targetRepository) CompleteTarget(id uint) error {
	return r.db.Model(&models.Target{}).Where("id = ?", id).Update("is_completed", true).Error
}

func (r *targetRepository) UpdateNotes(id uint, notes string) error {
	return r.db.Model(&models.Target{}).Where("id = ?", id).Update("notes", notes).Error
}

func (r *targetRepository) CountByMissionID(missionID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Target{}).Where("mission_id = ?", missionID).Count(&count).Error
	return count, err
}
