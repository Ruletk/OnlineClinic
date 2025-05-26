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

type Specialization struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"unique;not null"` // например, "Кардиолог"
	Description string    `gorm:"default:null"`
}

type ScheduleSlot struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	DoctorID      uuid.UUID  `gorm:"type:uuid;not null"`
	Date          time.Time  `gorm:"type:date;not null"`
	StartTime     time.Time  `gorm:"type:time;not null"` // например, "09:00:00"
	EndTime       time.Time  `gorm:"type:time;not null"` // например, "09:30:00"
	IsAvailable   bool       `gorm:"default:true"`
	AppointmentID *uuid.UUID `gorm:"type:uuid;default:null"` // ссылка на внешний сервис (nullable)
	MeetingLink   string     `gorm:"default:null"`           // для онлайн-консультаций
}
