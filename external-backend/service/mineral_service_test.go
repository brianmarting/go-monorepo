package service

import (
	"common/model"
	tracing_mocks "common/observability/tracing/mocks"
	"context"
	"errors"
	"external_backend/queue/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mineral = model.Mineral{}

func Test_mineralService_AddMineral(t *testing.T) {
	publisherMock := new(mocks.MineralPublisherMock)
	service := mineralService{
		tracer:    tracing_mocks.NewTracerMock(),
		publisher: publisherMock,
	}

	tests := []struct {
		name        string
		mockFn      func() *mock.Call
		expectError bool
		error       string
	}{
		{
			name: "Should add mineral",
			mockFn: func() *mock.Call {
				return publisherMock.On("Publish", mock.Anything, mineral).Return(nil)
			},
		},
		{
			name: "Should get err when adding mineral",
			mockFn: func() *mock.Call {
				return publisherMock.On("Publish", mock.Anything, mineral).Return(errors.New("failed to publish"))
			},
			expectError: true,
			error:       "failed to publish",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockCall *mock.Call
			if tt.mockFn != nil {
				mockCall = tt.mockFn()
			}
			defer func() {
				if mockCall != nil {
					mockCall.Unset()
				}
			}()
			err := service.AddMineral(context.Background(), mineral)

			if tt.expectError {
				assert.Equal(t, errors.New(tt.error), err)
			}

			publisherMock.AssertExpectations(t)
		})
	}
}
