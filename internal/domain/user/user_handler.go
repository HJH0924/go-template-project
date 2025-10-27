// Package user provides HTTP handlers for user service.
package user

import (
	"context"

	"github.com/HJH0924/go-template-project/internal/domain/user/service"
	userv1 "github.com/HJH0924/go-template-project/sdk/go/user/v1"
	"github.com/HJH0924/go-template-project/sdk/go/user/v1/userv1connect"

	"connectrpc.com/connect"
)

// Handler 用户服务handler.
type Handler struct {
	service service.UserService
}

// 确保实现userv1connect.UserServiceHandler接口.
var _ userv1connect.UserServiceHandler = (*Handler)(nil)

// NewHandler 创建handler实例.
func NewHandler(svc service.UserService) userv1connect.UserServiceHandler {
	return &Handler{
		service: svc,
	}
}

// SayHello 实现SayHello接口.
func (h *Handler) SayHello(
	ctx context.Context,
	req *connect.Request[userv1.SayHelloRequest],
) (*connect.Response[userv1.SayHelloResponse], error) {
	// 调用service层
	message, err := h.service.SayHello(ctx, req.Msg.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 构造响应
	resp := &userv1.SayHelloResponse{
		Message: message,
	}

	return connect.NewResponse(resp), nil
}
