package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spy-cat-agency/internal/models"
	"spy-cat-agency/internal/repository"
	"time"
)

type CatService interface {
	CreateCat(req *models.CreateCatRequest) (*models.SpyCat, error)
	GetCat(id uint) (*models.SpyCat, error)
	ListCats() ([]models.SpyCat, error)
	UpdateCat(id uint, req *models.UpdateCatRequest) (*models.SpyCat, error)
	DeleteCat(id uint) error
	ValidateBreed(breed string) error
}

type catService struct {
	catRepo repository.CatRepository
}

func NewCatService(catRepo repository.CatRepository) CatService {
	return &catService{catRepo: catRepo}
}

func (s *catService) CreateCat(req *models.CreateCatRequest) (*models.SpyCat, error) {
	// Validate breed
	if err := s.ValidateBreed(req.Breed); err != nil {
		return nil, fmt.Errorf("invalid breed: %w", err)
	}

	cat := &models.SpyCat{
		Name:            req.Name,
		YearsExperience: req.YearsExperience,
		Breed:           req.Breed,
		Salary:          req.Salary,
		IsAvailable:     true,
	}

	if err := s.catRepo.Create(cat); err != nil {
		return nil, fmt.Errorf("failed to create cat: %w", err)
	}

	return cat, nil
}

func (s *catService) GetCat(id uint) (*models.SpyCat, error) {
	return s.catRepo.GetByID(id)
}

func (s *catService) ListCats() ([]models.SpyCat, error) {
	return s.catRepo.GetAll()
}

func (s *catService) UpdateCat(id uint, req *models.UpdateCatRequest) (*models.SpyCat, error) {
	cat, err := s.catRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("cat not found: %w", err)
	}

	cat.Salary = req.Salary

	if err := s.catRepo.Update(cat); err != nil {
		return nil, fmt.Errorf("failed to update cat: %w", err)
	}

	return cat, nil
}

func (s *catService) DeleteCat(id uint) error {
	// Check if cat has active mission
	// This will be implemented when we have mission service
	return s.catRepo.Delete(id)
}

func (s *catService) ValidateBreed(breed string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch breeds: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var breeds []struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(body, &breeds); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, b := range breeds {
		if b.Name == breed {
			return nil
		}
	}

	return fmt.Errorf("breed '%s' not found in TheCatAPI", breed)
}
