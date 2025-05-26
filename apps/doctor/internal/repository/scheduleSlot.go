package repository

import (
	"doctor/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleSlotRepository interface {
	Create(slot *model.ScheduleSlot) error
	GetByID(id uuid.UUID) (*model.ScheduleSlot, error)
	Update(slot *model.ScheduleSlot) error
	Delete(id uuid.UUID) error
}

type scheduleslotRepo struct {
	db *gorm.DB
}

func NewScheduleSlotRepository(db *gorm.DB) ScheduleSlotRepository {
	return &scheduleslotRepo{db: db}
}

func (r *scheduleslotRepo) Create(slot *model.ScheduleSlot) error {
	return r.db.Create(slot).Error
}

func (r *scheduleslotRepo) GetByID(id uuid.UUID) (*model.ScheduleSlot, error) {
	var slot model.ScheduleSlot
	if err := r.db.First(&slot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *scheduleslotRepo) Update(slot *model.ScheduleSlot) error {
	return r.db.Save(slot).Error
}

func (r *scheduleslotRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ScheduleSlot{}, "id = ?", id).Error
}
