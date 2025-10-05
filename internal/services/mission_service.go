package services

import (
	"fmt"
	"spy-cat-agency/internal/models"
	"spy-cat-agency/internal/repository"
)

type MissionService interface {
	CreateMission(req *models.CreateMissionRequest) (*models.Mission, error)
	GetMission(id uint) (*models.Mission, error)
	ListMissions() ([]models.Mission, error)
	UpdateMission(id uint, req *models.UpdateMissionRequest) (*models.Mission, error)
	DeleteMission(id uint) error
	AssignCat(missionID, catID uint) error
	CompleteMission(missionID uint) error
	AddTarget(missionID uint, req *models.AddTargetRequest) (*models.Target, error)
	UpdateTarget(missionID, targetID uint, req *models.UpdateTargetRequest) (*models.Target, error)
	DeleteTarget(missionID, targetID uint) error
	CompleteTarget(missionID, targetID uint) error
	UpdateTargetNotes(missionID, targetID uint, req *models.UpdateTargetNotesRequest) error
}

type missionService struct {
	missionRepo repository.MissionRepository
	targetRepo  repository.TargetRepository
	catRepo     repository.CatRepository
}

func NewMissionService(missionRepo repository.MissionRepository, targetRepo repository.TargetRepository, catRepo repository.CatRepository) MissionService {
	return &missionService{
		missionRepo: missionRepo,
		targetRepo:  targetRepo,
		catRepo:     catRepo,
	}
}

func (s *missionService) CreateMission(req *models.CreateMissionRequest) (*models.Mission, error) {
	mission := &models.Mission{
		CatID:       req.CatID,
		IsCompleted: false,
	}

	if req.CatID != nil {
		cat, err := s.catRepo.GetByID(*req.CatID)
		if err != nil {
			return nil, fmt.Errorf("cat not found: %w", err)
		}
		if !cat.IsAvailable {
			return nil, fmt.Errorf("cat is not available")
		}
	}

	if err := s.missionRepo.Create(mission); err != nil {
		return nil, fmt.Errorf("failed to create mission: %w", err)
	}

	for _, targetReq := range req.Targets {
		target := &models.Target{
			MissionID:   mission.ID,
			Name:        targetReq.Name,
			Country:     targetReq.Country,
			IsCompleted: false,
		}
		if err := s.targetRepo.Create(target); err != nil {
			return nil, fmt.Errorf("failed to create target: %w", err)
		}
		mission.Targets = append(mission.Targets, *target)
	}

	if req.CatID != nil {
		if err := s.catRepo.SetAvailability(*req.CatID, false); err != nil {
			return nil, fmt.Errorf("failed to update cat availability: %w", err)
		}
	}

	return mission, nil
}

func (s *missionService) GetMission(id uint) (*models.Mission, error) {
	return s.missionRepo.GetByID(id)
}

func (s *missionService) ListMissions() ([]models.Mission, error) {
	return s.missionRepo.GetAll()
}

func (s *missionService) UpdateMission(id uint, req *models.UpdateMissionRequest) (*models.Mission, error) {
	mission, err := s.missionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted {
		return nil, fmt.Errorf("cannot update completed mission")
	}

	if req.CatID != nil && (mission.CatID == nil || *mission.CatID != *req.CatID) {
		if mission.CatID != nil {
			if err := s.catRepo.SetAvailability(*mission.CatID, true); err != nil {
				return nil, fmt.Errorf("failed to free up current cat: %w", err)
			}
		}

		cat, err := s.catRepo.GetByID(*req.CatID)
		if err != nil {
			return nil, fmt.Errorf("cat not found: %w", err)
		}
		if !cat.IsAvailable {
			return nil, fmt.Errorf("cat is not available")
		}

		mission.CatID = req.CatID
		if err := s.catRepo.SetAvailability(*req.CatID, false); err != nil {
			return nil, fmt.Errorf("failed to assign new cat: %w", err)
		}
	}

	if err := s.missionRepo.Update(mission); err != nil {
		return nil, fmt.Errorf("failed to update mission: %w", err)
	}

	return mission, nil
}

func (s *missionService) DeleteMission(id uint) error {
	mission, err := s.missionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}

	if mission.CatID != nil {
		return fmt.Errorf("cannot delete mission that is assigned to a cat")
	}

	return s.missionRepo.Delete(id)
}

func (s *missionService) AssignCat(missionID, catID uint) error {
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted {
		return fmt.Errorf("cannot assign cat to completed mission")
	}

	cat, err := s.catRepo.GetByID(catID)
	if err != nil {
		return fmt.Errorf("cat not found: %w", err)
	}

	if !cat.IsAvailable {
		return fmt.Errorf("cat is not available")
	}

	if mission.CatID != nil {
		if err := s.catRepo.SetAvailability(*mission.CatID, true); err != nil {
			return fmt.Errorf("failed to free up current cat: %w", err)
		}
	}

	if err := s.missionRepo.AssignCat(missionID, catID); err != nil {
		return fmt.Errorf("failed to assign cat: %w", err)
	}

	if err := s.catRepo.SetAvailability(catID, false); err != nil {
		return fmt.Errorf("failed to update cat availability: %w", err)
	}

	return nil
}

func (s *missionService) CompleteMission(missionID uint) error {
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted {
		return fmt.Errorf("mission is already completed")
	}

	targets, err := s.targetRepo.GetByMissionID(missionID)
	if err != nil {
		return fmt.Errorf("failed to get targets: %w", err)
	}

	if len(targets) == 0 {
		return fmt.Errorf("mission has no targets")
	}

	for _, target := range targets {
		if !target.IsCompleted {
			return fmt.Errorf("cannot complete mission: not all targets are completed")
		}
	}

	if err := s.missionRepo.CompleteMission(missionID); err != nil {
		return fmt.Errorf("failed to complete mission: %w", err)
	}

	if mission.CatID != nil {
		if err := s.catRepo.SetAvailability(*mission.CatID, true); err != nil {
			return fmt.Errorf("failed to free up cat: %w", err)
		}
	}

	return nil
}

func (s *missionService) AddTarget(missionID uint, req *models.AddTargetRequest) (*models.Target, error) {
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return nil, fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted {
		return nil, fmt.Errorf("cannot add target to completed mission")
	}

	count, err := s.targetRepo.CountByMissionID(missionID)
	if err != nil {
		return nil, fmt.Errorf("failed to count targets: %w", err)
	}

	if count >= 3 {
		return nil, fmt.Errorf("mission already has maximum number of targets (3)")
	}

	target := &models.Target{
		MissionID:   missionID,
		Name:        req.Name,
		Country:     req.Country,
		IsCompleted: false,
	}

	if err := s.targetRepo.Create(target); err != nil {
		return nil, fmt.Errorf("failed to create target: %w", err)
	}

	return target, nil
}

func (s *missionService) UpdateTarget(missionID, targetID uint, req *models.UpdateTargetRequest) (*models.Target, error) {
	target, err := s.targetRepo.GetByID(targetID)
	if err != nil {
		return nil, fmt.Errorf("target not found: %w", err)
	}

	if target.MissionID != missionID {
		return nil, fmt.Errorf("target does not belong to this mission")
	}

	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return nil, fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted || target.IsCompleted {
		return nil, fmt.Errorf("cannot update target in completed mission or completed target")
	}

	target.Name = req.Name
	target.Country = req.Country

	if err := s.targetRepo.Update(target); err != nil {
		return nil, fmt.Errorf("failed to update target: %w", err)
	}

	return target, nil
}

func (s *missionService) DeleteTarget(missionID, targetID uint) error {
	target, err := s.targetRepo.GetByID(targetID)
	if err != nil {
		return fmt.Errorf("target not found: %w", err)
	}

	if target.MissionID != missionID {
		return fmt.Errorf("target does not belong to this mission")
	}

	if target.IsCompleted {
		return fmt.Errorf("cannot delete completed target")
	}

	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted {
		return fmt.Errorf("cannot delete target from completed mission")
	}

	return s.targetRepo.Delete(targetID)
}

func (s *missionService) CompleteTarget(missionID, targetID uint) error {
	target, err := s.targetRepo.GetByID(targetID)
	if err != nil {
		return fmt.Errorf("target not found: %w", err)
	}

	if target.MissionID != missionID {
		return fmt.Errorf("target does not belong to this mission")
	}

	if target.IsCompleted {
		return fmt.Errorf("target is already completed")
	}

	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted {
		return fmt.Errorf("cannot complete target in completed mission")
	}

	return s.targetRepo.CompleteTarget(targetID)
}

func (s *missionService) UpdateTargetNotes(missionID, targetID uint, req *models.UpdateTargetNotesRequest) error {
	target, err := s.targetRepo.GetByID(targetID)
	if err != nil {
		return fmt.Errorf("target not found: %w", err)
	}

	if target.MissionID != missionID {
		return fmt.Errorf("target does not belong to this mission")
	}

	if target.IsCompleted {
		return fmt.Errorf("cannot update notes for completed target")
	}

	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}

	if mission.IsCompleted {
		return fmt.Errorf("cannot update notes in completed mission")
	}

	newNotes := req.Notes
	if existing := target.Notes; len(existing) > 0 {
		newNotes = existing + "\n" + req.Notes
	}

	return s.targetRepo.UpdateNotes(targetID, newNotes)
}
