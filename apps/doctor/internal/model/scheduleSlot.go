package model

import (
	"github.com/google/uuid"
	"time"
)

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
