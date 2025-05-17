package model

type Doctor struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Specialty string `json:"specialty" db:"specialty"`
	Email     string `json:"email" db:"email"`
}
