// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package repository_mock

import (
	"auth/internal/repository"

	mock "github.com/stretchr/testify/mock"
)

// NewMockRoleRepository creates a new instance of MockRoleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRoleRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRoleRepository {
	mock := &MockRoleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockRoleRepository is an autogenerated mock type for the RoleRepository type
type MockRoleRepository struct {
	mock.Mock
}

type MockRoleRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRoleRepository) EXPECT() *MockRoleRepository_Expecter {
	return &MockRoleRepository_Expecter{mock: &_m.Mock}
}

// CreateRole provides a mock function for the type MockRoleRepository
func (_mock *MockRoleRepository) CreateRole(name string) error {
	ret := _mock.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for CreateRole")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string) error); ok {
		r0 = returnFunc(name)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRoleRepository_CreateRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRole'
type MockRoleRepository_CreateRole_Call struct {
	*mock.Call
}

// CreateRole is a helper method to define mock.On call
//   - name
func (_e *MockRoleRepository_Expecter) CreateRole(name interface{}) *MockRoleRepository_CreateRole_Call {
	return &MockRoleRepository_CreateRole_Call{Call: _e.mock.On("CreateRole", name)}
}

func (_c *MockRoleRepository_CreateRole_Call) Run(run func(name string)) *MockRoleRepository_CreateRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRoleRepository_CreateRole_Call) Return(err error) *MockRoleRepository_CreateRole_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRoleRepository_CreateRole_Call) RunAndReturn(run func(name string) error) *MockRoleRepository_CreateRole_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteRole provides a mock function for the type MockRoleRepository
func (_mock *MockRoleRepository) DeleteRole(name string) error {
	ret := _mock.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRole")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string) error); ok {
		r0 = returnFunc(name)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRoleRepository_DeleteRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteRole'
type MockRoleRepository_DeleteRole_Call struct {
	*mock.Call
}

// DeleteRole is a helper method to define mock.On call
//   - name
func (_e *MockRoleRepository_Expecter) DeleteRole(name interface{}) *MockRoleRepository_DeleteRole_Call {
	return &MockRoleRepository_DeleteRole_Call{Call: _e.mock.On("DeleteRole", name)}
}

func (_c *MockRoleRepository_DeleteRole_Call) Run(run func(name string)) *MockRoleRepository_DeleteRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRoleRepository_DeleteRole_Call) Return(err error) *MockRoleRepository_DeleteRole_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRoleRepository_DeleteRole_Call) RunAndReturn(run func(name string) error) *MockRoleRepository_DeleteRole_Call {
	_c.Call.Return(run)
	return _c
}

// GetRoleByName provides a mock function for the type MockRoleRepository
func (_mock *MockRoleRepository) GetRoleByName(name string) (*repository.Role, error) {
	ret := _mock.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetRoleByName")
	}

	var r0 *repository.Role
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (*repository.Role, error)); ok {
		return returnFunc(name)
	}
	if returnFunc, ok := ret.Get(0).(func(string) *repository.Role); ok {
		r0 = returnFunc(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Role)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(name)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRoleRepository_GetRoleByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRoleByName'
type MockRoleRepository_GetRoleByName_Call struct {
	*mock.Call
}

// GetRoleByName is a helper method to define mock.On call
//   - name
func (_e *MockRoleRepository_Expecter) GetRoleByName(name interface{}) *MockRoleRepository_GetRoleByName_Call {
	return &MockRoleRepository_GetRoleByName_Call{Call: _e.mock.On("GetRoleByName", name)}
}

func (_c *MockRoleRepository_GetRoleByName_Call) Run(run func(name string)) *MockRoleRepository_GetRoleByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRoleRepository_GetRoleByName_Call) Return(role *repository.Role, err error) *MockRoleRepository_GetRoleByName_Call {
	_c.Call.Return(role, err)
	return _c
}

func (_c *MockRoleRepository_GetRoleByName_Call) RunAndReturn(run func(name string) (*repository.Role, error)) *MockRoleRepository_GetRoleByName_Call {
	_c.Call.Return(run)
	return _c
}
