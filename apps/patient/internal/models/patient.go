package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	Name          string
	Age           int
	MedicalRecord string `gorm:"type:text"`
	BloodType     string
	Allergies     string `gorm:"type:text"`
	Medications   string `gorm:"type:text"`
}
