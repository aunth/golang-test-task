package repository

import (
	"spy-cat-agency/internal/models"

	"gorm.io/gorm"
)

type MissionRepository interface {
	Create(mission *models.Mission) error
	GetByID(id uint) (*models.Mission, error)
	GetAll() ([]models.Mission, error)
	Update(mission *models.Mission) error
	Delete(id uint) error
	GetByCatID(catID uint) (*models.Mission, error)
	AssignCat(missionID, catID uint) error
	CompleteMission(missionID uint) error
}

type missionRepository struct {
	db *gorm.DB
}

func NewMissionRepository(db *gorm.DB) MissionRepository {
	return &missionRepository{db: db}
}

func (r *missionRepository) Create(mission *models.Mission) error {
	return r.db.Create(mission).Error
}

func (r *missionRepository) GetByID(id uint) (*models.Mission, error) {
	var mission models.Mission
	err := r.db.Preload("Cat").Preload("Targets").First(&mission, id).Error
	if err != nil {
		return nil, err
	}
	return &mission, nil
}

func (r *missionRepository) GetAll() ([]models.Mission, error) {
	var missions []models.Mission
	err := r.db.Preload("Cat").Preload("Targets").Find(&missions).Error
	return missions, err
}

func (r *missionRepository) Update(mission *models.Mission) error {
	return r.db.Save(mission).Error
}

func (r *missionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Mission{}, id).Error
}

func (r *missionRepository) GetByCatID(catID uint) (*models.Mission, error) {
	var mission models.Mission
	err := r.db.Preload("Cat").Preload("Targets").Where("cat_id = ? AND is_completed = ?", catID, false).First(&mission).Error
	if err != nil {
		return nil, err
	}
	return &mission, nil
}

func (r *missionRepository) AssignCat(missionID, catID uint) error {
	return r.db.Model(&models.Mission{}).Where("id = ?", missionID).Update("cat_id", catID).Error
}

func (r *missionRepository) CompleteMission(missionID uint) error {
	return r.db.Model(&models.Mission{}).Where("id = ?", missionID).Update("is_completed", true).Error
}
