// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Database is an autogenerated mock type for the Database type
type Database struct {
	mock.Mock
}

// GetURL provides a mock function with given fields: id
func (_m *Database) GetURL(id uint64) (string, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (string, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint64) string); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetURLs provides a mock function with given fields: userID
func (_m *Database) GetURLs(userID int) ([]URL, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetURLs")
	}

	var r0 []URL
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]URL, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int) []URL); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]URL)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: username
func (_m *Database) GetUser(username string) (*User, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*User, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) *User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertURL provides a mock function with given fields: userID, url
func (_m *Database) InsertURL(userID int, url string) (int64, string, error) {
	ret := _m.Called(userID, url)

	if len(ret) == 0 {
		panic("no return value specified for InsertURL")
	}

	var r0 int64
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(int, string) (int64, string, error)); ok {
		return rf(userID, url)
	}
	if rf, ok := ret.Get(0).(func(int, string) int64); ok {
		r0 = rf(userID, url)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int, string) string); ok {
		r1 = rf(userID, url)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(int, string) error); ok {
		r2 = rf(userID, url)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Ping provides a mock function with given fields:
func (_m *Database) Ping() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ping")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SignupUser provides a mock function with given fields: username, password, email
func (_m *Database) SignupUser(username string, password string, email string) error {
	ret := _m.Called(username, password, email)

	if len(ret) == 0 {
		panic("no return value specified for SignupUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(username, password, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateUser provides a mock function with given fields: username
func (_m *Database) ValidateUser(username string) error {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for ValidateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDatabase creates a new instance of Database. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *Database {
	mock := &Database{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
