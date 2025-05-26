package model

import (
	"github.com/google/uuid"
	"time"
)

type DoctorStatus string

const (
	Active   DoctorStatus = "ACTIVE"
	OnLeave  DoctorStatus = "ON_LEAVE"
	Inactive DoctorStatus = "INACTIVE"
)

type Doctor struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey"`
	FirstName        string         `gorm:"not null"`
	LastName         string         `gorm:"not null"`
	Patronymic       string         `gorm:"default:null"` // опционально
	DateOfBirth      time.Time      `gorm:"type:date"`
	Specialization   Specialization `gorm:"foreignKey:SpecializationID"`
	SpecializationID uuid.UUID      `gorm:"type:uuid"`
	Status           DoctorStatus   `gorm:"type:varchar(20);default:'ACTIVE'"`
	ScheduleSlots    []ScheduleSlot `gorm:"foreignKey:DoctorID"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
}
