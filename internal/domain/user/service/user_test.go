package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_SayHello(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with name",
			input:    "Alice",
			expected: "Hello, Alice!",
		},
		{
			name:     "empty name defaults to World",
			input:    "",
			expected: "Hello, World!",
		},
		{
			name:     "chinese name",
			input:    "小明",
			expected: "Hello, 小明!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			service := NewUserService()
			resp, err := service.SayHello(ctx, tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, resp)
		})
	}
}
