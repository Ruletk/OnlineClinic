package repository

import "gorm.io/gorm"

type Role struct {
	ID   int64  `json:"id" gorm:"primaryKey;column:id"`
	Name string `json:"name" gorm:"unique"`
}

func (Role) TableName() string {
	return "roles"
}

func GetRoleNames(roles []Role) []string {
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames
}

type RoleRepository interface {
	GetRoleByName(name string) (*Role, error)
	CreateRole(name string) error
	DeleteRole(name string) error
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(DB *gorm.DB) RoleRepository {
	return roleRepository{DB: DB}
}

func (r roleRepository) GetRoleByName(name string) (*Role, error) {
	var role Role
	err := r.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r roleRepository) CreateRole(name string) error {
	role := Role{Name: name}
	return r.DB.Create(&role).Error
}

func (r roleRepository) DeleteRole(name string) error {
	return r.DB.Where("name = ?", name).Delete(&Role{}).Error
}
