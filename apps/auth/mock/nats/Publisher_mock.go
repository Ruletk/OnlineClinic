// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package nats_mock

import (
	mock "github.com/stretchr/testify/mock"
)

// NewMockPublisher creates a new instance of MockPublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPublisher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPublisher {
	mock := &MockPublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockPublisher is an autogenerated mock type for the Publisher type
type MockPublisher struct {
	mock.Mock
}

type MockPublisher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPublisher) EXPECT() *MockPublisher_Expecter {
	return &MockPublisher_Expecter{mock: &_m.Mock}
}

// PublishEmailMessage provides a mock function for the type MockPublisher
func (_mock *MockPublisher) PublishEmailMessage(to string, subject string, message string) error {
	ret := _mock.Called(to, subject, message)

	if len(ret) == 0 {
		panic("no return value specified for PublishEmailMessage")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = returnFunc(to, subject, message)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockPublisher_PublishEmailMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishEmailMessage'
type MockPublisher_PublishEmailMessage_Call struct {
	*mock.Call
}

// PublishEmailMessage is a helper method to define mock.On call
//   - to
//   - subject
//   - message
func (_e *MockPublisher_Expecter) PublishEmailMessage(to interface{}, subject interface{}, message interface{}) *MockPublisher_PublishEmailMessage_Call {
	return &MockPublisher_PublishEmailMessage_Call{Call: _e.mock.On("PublishEmailMessage", to, subject, message)}
}

func (_c *MockPublisher_PublishEmailMessage_Call) Run(run func(to string, subject string, message string)) *MockPublisher_PublishEmailMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockPublisher_PublishEmailMessage_Call) Return(err error) *MockPublisher_PublishEmailMessage_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockPublisher_PublishEmailMessage_Call) RunAndReturn(run func(to string, subject string, message string) error) *MockPublisher_PublishEmailMessage_Call {
	_c.Call.Return(run)
	return _c
}
