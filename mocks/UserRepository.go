// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "drdgvhbh/discordbot/internal/user/entity"
import mock "github.com/stretchr/testify/mock"

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// InsertUser provides a mock function with given fields: user
func (_m *UserRepository) InsertUser(user *entity.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
