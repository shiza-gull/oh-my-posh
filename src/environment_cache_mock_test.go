package main

import mock "github.com/stretchr/testify/mock"

// MockedCache is an autogenerated mock type for the cache type
type MockedCache struct {
	mock.Mock
}

// close provides a mock function with given fields:
func (_m *MockedCache) close() {
	_m.Called()
}

// get provides a mock function with given fields: key
func (_m *MockedCache) get(key string) (string, bool) {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// init provides a mock function with given fields: home
func (_m *MockedCache) init(home string) {
	_m.Called(home)
}

// set provides a mock function with given fields: key, value, ttl
func (_m *MockedCache) set(key, value string, ttl int) {
	_m.Called(key, value, ttl)
}
