package mocks

import (
	"context"

	"go-api-server/internal/adapter/repository/postgresql/schema"

	"github.com/stretchr/testify/mock"
)

type MockConfigurationRepository struct {
	mock.Mock
}

func NewMockConfigurationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConfigurationRepository {
	mock := &MockConfigurationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

func (_m *MockConfigurationRepository) FindByKey(ctx context.Context, key string) (*schema.FindByKeySchema, error) {
	ret := _m.Called(ctx, key)

	var r0 *schema.FindByKeySchema
	if rf, ok := ret.Get(0).(func(context.Context, string) *schema.FindByKeySchema); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*schema.FindByKeySchema)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
