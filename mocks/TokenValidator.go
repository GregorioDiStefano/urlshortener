// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	mock "github.com/stretchr/testify/mock"
)

// TokenValidator is an autogenerated mock type for the TokenValidator type
type TokenValidator struct {
	mock.Mock
}

// GenerateUserToken provides a mock function with given fields: userid
func (_m *TokenValidator) GenerateUserToken(userid int) (string, error) {
	ret := _m.Called(userid)

	if len(ret) == 0 {
		panic("no return value specified for GenerateUserToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (string, error)); ok {
		return rf(userid)
	}
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(userid)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// JWTMiddleware provides a mock function with given fields:
func (_m *TokenValidator) JWTMiddleware() gin.HandlerFunc {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for JWTMiddleware")
	}

	var r0 gin.HandlerFunc
	if rf, ok := ret.Get(0).(func() gin.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.HandlerFunc)
		}
	}

	return r0
}

// ValidateToken provides a mock function with given fields: token
func (_m *TokenValidator) ValidateToken(token string) (*jwt.Token, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *jwt.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*jwt.Token, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) *jwt.Token); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTokenValidator creates a new instance of TokenValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenValidator {
	mock := &TokenValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
