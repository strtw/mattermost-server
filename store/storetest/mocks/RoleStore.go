// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost-server/v5/model"
	mock "github.com/stretchr/testify/mock"
)

// RoleStore is an autogenerated mock type for the RoleStore type
type RoleStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: roleId
func (_m *RoleStore) Delete(roleId string) (*model.Role, *model.AppError) {
	ret := _m.Called(roleId)

	var r0 *model.Role
	if rf, ok := ret.Get(0).(func(string) *model.Role); ok {
		r0 = rf(roleId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Role)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(roleId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// Get provides a mock function with given fields: roleId
func (_m *RoleStore) Get(roleId string) (*model.Role, *model.AppError) {
	ret := _m.Called(roleId)

	var r0 *model.Role
	if rf, ok := ret.Get(0).(func(string) *model.Role); ok {
		r0 = rf(roleId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Role)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(roleId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *RoleStore) GetAll() ([]*model.Role, *model.AppError) {
	ret := _m.Called()

	var r0 []*model.Role
	if rf, ok := ret.Get(0).(func() []*model.Role); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Role)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func() *model.AppError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *RoleStore) GetByName(name string) (*model.Role, *model.AppError) {
	ret := _m.Called(name)

	var r0 *model.Role
	if rf, ok := ret.Get(0).(func(string) *model.Role); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Role)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetByNames provides a mock function with given fields: names
func (_m *RoleStore) GetByNames(names []string) ([]*model.Role, *model.AppError) {
	ret := _m.Called(names)

	var r0 []*model.Role
	if rf, ok := ret.Get(0).(func([]string) []*model.Role); ok {
		r0 = rf(names)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Role)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func([]string) *model.AppError); ok {
		r1 = rf(names)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// HigherScopedPermissions provides a mock function with given fields: roleNames
func (_m *RoleStore) HigherScopedPermissions(roleNames []string) ([]*model.RolePermissions, *model.AppError) {
	ret := _m.Called(roleNames)

	var r0 []*model.RolePermissions
	if rf, ok := ret.Get(0).(func([]string) []*model.RolePermissions); ok {
		r0 = rf(roleNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.RolePermissions)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func([]string) *model.AppError); ok {
		r1 = rf(roleNames)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// PermanentDeleteAll provides a mock function with given fields:
func (_m *RoleStore) PermanentDeleteAll() *model.AppError {
	ret := _m.Called()

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func() *model.AppError); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// Save provides a mock function with given fields: role
func (_m *RoleStore) Save(role *model.Role) (*model.Role, *model.AppError) {
	ret := _m.Called(role)

	var r0 *model.Role
	if rf, ok := ret.Get(0).(func(*model.Role) *model.Role); ok {
		r0 = rf(role)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Role)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Role) *model.AppError); ok {
		r1 = rf(role)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}
