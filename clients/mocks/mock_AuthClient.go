// Code generated by mockery v2.33.2. DO NOT EDIT.

package apiclientsmocks

import (
	context "context"

	apiclients "github.com/a-novel/go-apis/clients"

	mock "github.com/stretchr/testify/mock"
)

// AuthClient is an autogenerated mock type for the AuthClient type
type AuthClient struct {
	mock.Mock
}

type AuthClient_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthClient) EXPECT() *AuthClient_Expecter {
	return &AuthClient_Expecter{mock: &_m.Mock}
}

// IntrospectToken provides a mock function with given fields: ctx, token
func (_m *AuthClient) IntrospectToken(ctx context.Context, token string) (*apiclients.UserTokenStatus, error) {
	ret := _m.Called(ctx, token)

	var r0 *apiclients.UserTokenStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*apiclients.UserTokenStatus, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *apiclients.UserTokenStatus); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apiclients.UserTokenStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthClient_IntrospectToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IntrospectToken'
type AuthClient_IntrospectToken_Call struct {
	*mock.Call
}

// IntrospectToken is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *AuthClient_Expecter) IntrospectToken(ctx interface{}, token interface{}) *AuthClient_IntrospectToken_Call {
	return &AuthClient_IntrospectToken_Call{Call: _e.mock.On("IntrospectToken", ctx, token)}
}

func (_c *AuthClient_IntrospectToken_Call) Run(run func(ctx context.Context, token string)) *AuthClient_IntrospectToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *AuthClient_IntrospectToken_Call) Return(_a0 *apiclients.UserTokenStatus, _a1 error) *AuthClient_IntrospectToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthClient_IntrospectToken_Call) RunAndReturn(run func(context.Context, string) (*apiclients.UserTokenStatus, error)) *AuthClient_IntrospectToken_Call {
	_c.Call.Return(run)
	return _c
}

// Ping provides a mock function with given fields: ctx
func (_m *AuthClient) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthClient_Ping_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ping'
type AuthClient_Ping_Call struct {
	*mock.Call
}

// Ping is a helper method to define mock.On call
//   - ctx context.Context
func (_e *AuthClient_Expecter) Ping(ctx interface{}) *AuthClient_Ping_Call {
	return &AuthClient_Ping_Call{Call: _e.mock.On("Ping", ctx)}
}

func (_c *AuthClient_Ping_Call) Run(run func(ctx context.Context)) *AuthClient_Ping_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *AuthClient_Ping_Call) Return(_a0 error) *AuthClient_Ping_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthClient_Ping_Call) RunAndReturn(run func(context.Context) error) *AuthClient_Ping_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuthClient creates a new instance of AuthClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthClient {
	mock := &AuthClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
