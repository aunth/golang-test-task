package repository

import (
	"spy-cat-agency/internal/models"

	"gorm.io/gorm"
)

type CatRepository interface {
	Create(cat *models.SpyCat) error
	GetByID(id uint) (*models.SpyCat, error)
	GetAll() ([]models.SpyCat, error)
	Update(cat *models.SpyCat) error
	Delete(id uint) error
	GetAvailable() ([]models.SpyCat, error)
	SetAvailability(id uint, available bool) error
}

type catRepository struct {
	db *gorm.DB
}

func NewCatRepository(db *gorm.DB) CatRepository {
	return &catRepository{db: db}
}

func (r *catRepository) Create(cat *models.SpyCat) error {
	return r.db.Create(cat).Error
}

func (r *catRepository) GetByID(id uint) (*models.SpyCat, error) {
	var cat models.SpyCat
	err := r.db.First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *catRepository) GetAll() ([]models.SpyCat, error) {
	var cats []models.SpyCat
	err := r.db.Find(&cats).Error
	return cats, err
}

func (r *catRepository) Update(cat *models.SpyCat) error {
	return r.db.Save(cat).Error
}

func (r *catRepository) Delete(id uint) error {
	return r.db.Delete(&models.SpyCat{}, id).Error
}

func (r *catRepository) GetAvailable() ([]models.SpyCat, error) {
	var cats []models.SpyCat
	err := r.db.Where("is_available = ?", true).Find(&cats).Error
	return cats, err
}

func (r *catRepository) SetAvailability(id uint, available bool) error {
	return r.db.Model(&models.SpyCat{}).Where("id = ?", id).Update("is_available", available).Error
}
