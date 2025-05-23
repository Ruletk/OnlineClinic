package model

// Doctor — модель доктора, совпадает с таблицей doctors
type Doctor struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Name      string `gorm:"not null" json:"name"`
    Specialty string `gorm:"not null" json:"specialty"`
    Email     string `gorm:"not null;unique" json:"email"`
}
