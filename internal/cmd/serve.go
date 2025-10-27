package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HJH0924/go-template-project/internal/config"
	"github.com/HJH0924/go-template-project/internal/domain/user"
	"github.com/HJH0924/go-template-project/internal/domain/user/service"
	"github.com/HJH0924/go-template-project/sdk/go/user/v1/userv1connect"

	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server with Connect/gRPC support`,
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(_ *cobra.Command, _ []string) error {
	// 初始化logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting server")

	// 加载配置
	if err := config.Load(cfgFile); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg := config.Get()
	logger.Info("config loaded",
		slog.String("host", cfg.Server.Host),
		slog.Int("port", cfg.Server.Port),
	)

	// 初始化service和handler
	userService := service.NewUserService()
	userHandler := user.NewHandler(userService)

	// 创建HTTP服务器
	mux := http.NewServeMux()

	// 注册Connect服务
	path, handler := userv1connect.NewUserServiceHandler(userHandler)
	mux.Handle(path, handler)

	// 添加健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("OK")); err != nil {
			logger.Error("failed to write health check response", slog.Any("error", err))
		}
	})

	// 使用h2c支持HTTP/2
	server := &http.Server{
		Addr:              cfg.Server.Address(),
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Info("server listening", slog.String("addr", server.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("server stopped")

	return nil
}
