package queue

import (
	"common/model"
	tracing_mocks "common/observability/tracing/mocks"
	"common/queue/mocks"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mineral = model.Mineral{}

func Test_mineralPublisher_Publish(t *testing.T) {
	publisherMock := new(mocks.PublisherMock)
	publisher := mineralPublisher{
		tracer:    tracing_mocks.NewTracerMock(),
		publisher: publisherMock,
	}

	data, _ := json.Marshal(mineral)

	tests := []struct {
		name        string
		mockFn      func() *mock.Call
		expectError bool
		error       string
	}{
		{
			name: "Should publish msg",
			mockFn: func() *mock.Call {
				return publisherMock.On("Publish", mock.Anything, "mineral.deposit", data).Return(nil)
			},
		},
		{
			name: "Should fail on marshal err",
			mockFn: func() *mock.Call {
				return publisherMock.On("Publish", mock.Anything, "mineral.deposit", data).Return(errors.New("failed to marshal"))
			},
			expectError: true,
			error:       "failed to marshal",
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

			err := publisher.Publish(context.Background(), mineral)
			if tt.expectError {
				assert.Equal(t, errors.New(tt.error), err)
			}

			publisherMock.AssertExpectations(t)
		})
	}
}
