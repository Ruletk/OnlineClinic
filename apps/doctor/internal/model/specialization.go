package model

import "github.com/google/uuid"

type Specialization struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"unique;not null"` // например, "Кардиолог"
	Description string    `gorm:"default:null"`
}
