package service

import (
	"auth/internal/repository"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
)

type RoleService interface {
	Create(name string) error
	Delete(name string) error
	Get(name string) (*repository.Role, error)
}

type roleService struct {
	roleRepository repository.RoleRepository
}

func NewRoleService(roleRepository repository.RoleRepository) RoleService {
	return &roleService{roleRepository: roleRepository}
}

func (r *roleService) Create(name string) error {
	logging.Logger.Info("Creating new role: ", name)
	return r.roleRepository.CreateRole(name)
}

func (r *roleService) Delete(name string) error {
	logging.Logger.Info("Deleting role: ", name)
	return r.roleRepository.DeleteRole(name)
}

func (r *roleService) Get(name string) (*repository.Role, error) {
	logging.Logger.Info("Retrieving role: ", name)
	return r.roleRepository.GetRoleByName(name)
}
