// Package service provides business logic for user domain.
package service

import (
	"context"
	"fmt"
	"log/slog"
)

// UserService 用户服务接口.
type UserService interface {
	SayHello(ctx context.Context, name string) (string, error)
}

// userService 用户服务实现.
type userService struct{}

// NewUserService 创建用户服务实例.
func NewUserService() UserService {
	return &userService{}
}

// SayHello 打招呼.
func (s *userService) SayHello(ctx context.Context, name string) (string, error) {
	if name == "" {
		name = "World"
	}

	message := fmt.Sprintf("Hello, %s!", name)

	slog.InfoContext(ctx, "say hello called",
		slog.String("name", name),
		slog.String("message", message),
	)

	return message, nil
}
