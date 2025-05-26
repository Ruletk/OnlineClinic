package service

import (
	"doctor/internal/model"
	"doctor/internal/repository"
	"errors"

	"github.com/google/uuid"
)

var ErrSlotNotFound = errors.New("schedule slot not found")

type ScheduleSlotService interface {
	CreateSlot(req model.ScheduleSlot) (*model.ScheduleSlot, error)
	GetSlotByID(id uuid.UUID) (*model.ScheduleSlot, error)
	UpdateSlot(req model.ScheduleSlot) (*model.ScheduleSlot, error)
	DeleteSlot(id uuid.UUID) error
}

type scheduleslotService struct {
	repo repository.ScheduleSlotRepository
}

func NewScheduleSlotService(r repository.ScheduleSlotRepository) ScheduleSlotService {
	return &scheduleslotService{repo: r}
}

func (s *scheduleslotService) CreateSlot(req model.ScheduleSlot) (*model.ScheduleSlot, error) {
	req.ID = uuid.New()
	if err := s.repo.Create(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *scheduleslotService) GetSlotByID(id uuid.UUID) (*model.ScheduleSlot, error) {
	slot, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrSlotNotFound
	}
	return slot, nil
}

func (s *scheduleslotService) UpdateSlot(req model.ScheduleSlot) (*model.ScheduleSlot, error) {
	if _, err := s.repo.GetByID(req.ID); err != nil {
		return nil, ErrSlotNotFound
	}
	if err := s.repo.Update(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *scheduleslotService) DeleteSlot(id uuid.UUID) error {
	return s.repo.Delete(id)
}
