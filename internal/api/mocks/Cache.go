// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "gitlab.com/pet-pr-social-network/feed-service/internal/model"
)

// Cache is an autogenerated mock type for the Cache type
type Cache struct {
	mock.Mock
}

// GetPostsByUserID provides a mock function with given fields: ctx, userID
func (_m *Cache) GetPostsByUserID(ctx context.Context, userID int64) ([]model.Post, error) {
	ret := _m.Called(ctx, userID)

	var r0 []model.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]model.Post, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []model.Post); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCache interface {
	mock.TestingT
	Cleanup(func())
}

// NewCache creates a new instance of Cache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCache(t mockConstructorTestingTNewCache) *Cache {
	mock := &Cache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}