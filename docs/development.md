# Go Template Project 开发指南

## 目录

- [快速开始](#快速开始)
- [开发环境设置](#开发环境设置)
- [项目结构](#项目结构)
- [开发流程](#开发流程)
- [API 开发](#api-开发)
- [测试](#测试)
- [配置管理](#配置管理)
- [日志记录](#日志记录)
- [故障排除](#故障排除)

## 快速开始

### 前置要求

- **Go 1.24+**: [下载地址](https://go.dev/dl/)
- **Make**: 用于构建和运行任务
- **Git**: 版本控制
- **开发工具**（可通过 `make install-deps` 自动安装）:
  - golangci-lint: 代码检查工具
  - gofumpt: Go 代码格式化工具
  - goimports: Import 语句整理工具
  - goimports-reviser: Import 语句排序工具
  - buf: Protocol Buffers 工具

### 安装步骤

```bash
# 1. 克隆仓库
git clone https://github.com/HJH0924/go-template-project.git
cd go-template-project

# 2. 安装开发依赖（推荐）
make install-deps

# 3. 安装 Go 模块依赖
go mod tidy

# 4. 生成 proto 代码
make generate

# 5. 构建项目
make build

# 6. 运行服务
./bin/server serve --config configs/config.yaml
```

服务启动后将监听 `http://0.0.0.0:8080`。

### 验证安装

```bash
# 健康检查
curl http://localhost:8080/health

# 测试 API
curl --request POST \
  --url http://localhost:8080/user.v1.UserService/SayHello \
  --header 'Content-Type: application/json' \
  --data '{"name": "World"}'
```

## 开发环境设置

### 安装开发依赖

运行 `make install-deps` 会自动完成以下操作：

1. **检查并安装开发工具**：
   - golangci-lint: 用于代码质量检查
   - gofumpt: 比标准 gofmt 更严格的格式化工具
   - goimports: 自动添加、删除和整理 import 语句
   - goimports-reviser: 按规则排序 import 语句
   - buf: 用于 Protocol Buffers 的编译和管理

2. **安装 Git pre-commit hook**：
   - 自动在每次提交前运行代码格式化（`make format`）
   - 自动运行代码检查（`make lint`）
   - 如果格式化或检查失败，将阻止提交
   - 如果格式化产生了变更，会提示先添加这些变更

### Pre-commit Hook

Pre-commit hook 位于 `scripts/pre-commit`，会在每次 `git commit` 时自动执行。

**手动安装 hook：**
```bash
make install-hooks
```

**跳过 hook（不推荐）：**
```bash
git commit --no-verify -m "your message"
```

**Hook 工作流程：**
1. 运行 `make format` 格式化代码
2. 运行 `make lint` 进行代码检查
3. 检查是否有格式化产生的未提交更改
4. 如果所有检查通过，允许提交；否则阻止提交并提示修复

### 开发工具使用

```bash
# 格式化代码
make format

# 运行代码检查
make lint

# 构建项目（包含 format 和 lint）
make build

# 运行测试
make test

# 生成 proto 代码
make generate

# 开发模式运行（无需构建）
make dev
```

## 项目结构

```
go-template-project/
├── cmd/
│   └── server/              # 主程序入口
│       └── main.go
├── configs/
│   └── config.yaml          # 配置文件
├── docs/                    # 文档
│   ├── .vitepress/          # VitePress 配置
│   ├── development.md       # 本文档
│   └── index.md             # 文档首页
├── internal/                # 私有应用代码
│   ├── cmd/                 # Cobra 命令定义
│   │   ├── root.go          # 根命令
│   │   └── serve.go         # serve 子命令
│   ├── config/              # 配置管理
│   │   └── config.go        # Viper 配置加载
│   └── domain/              # 业务域
│       └── user/
│           ├── service/     # 业务逻辑层
│           │   ├── user.go  # Service 接口和实现
│           │   └── user_test.go  # 单元测试
│           └── user_handler.go   # Handler 层（Connect RPC）
├── proto/                   # Protocol Buffer 定义
│   └── user/
│       └── v1/
│           └── user.proto
├── sdk/                     # 生成的 SDK 代码（不要手动编辑）
│   └── go/
│       └── user/
│           └── v1/
├── scripts/                 # 构建脚本和工具
│   └── pre-commit           # Git pre-commit hook
├── buf.gen.yaml             # Buf 代码生成配置
├── buf.yaml                 # Buf 项目配置
├── Makefile                 # 构建任务
├── .golangci.toml           # golangci-lint 配置
└── go.mod
```

### 关键目录说明

- **cmd/**: 应用程序入口点
- **internal/**: 私有应用代码，不对外暴露
  - **cmd/**: Cobra 命令行工具定义
  - **config/**: 配置加载和管理
  - **domain/**: 领域驱动设计结构，每个域独立
- **proto/**: Protocol Buffer 定义（API 的真实来源）
- **sdk/**: 自动生成的代码（永远不要手动编辑）
- **scripts/**: 构建脚本和开发工具
- **docs/**: VitePress 文档

## 开发流程

### 1. 添加新 API

#### 步骤 1: 定义 Proto 接口

在 `proto/<service>/v1/` 中编辑或创建 `.proto` 文件：

```protobuf
syntax = "proto3";

package myservice.v1;

option go_package = "github.com/HJH0924/go-template-project/sdk/go/myservice/v1;myservicev1";

service MyService {
  rpc DoSomething(DoSomethingRequest) returns (DoSomethingResponse) {}
}

message DoSomethingRequest {
  string param = 1;
}

message DoSomethingResponse {
  string result = 1;
}
```

#### 步骤 2: 更新 Buf 配置

如果是新的服务目录，需要在 `buf.yaml` 中确认模块配置正确。

#### 步骤 3: 生成代码

```bash
make generate
```

这将生成：
- `sdk/go/myservice/v1/myservice.pb.go` - 消息定义
- `sdk/go/myservice/v1/myservicev1connect/myservice.connect.go` - 服务接口

#### 步骤 4: 实现 Service 层

创建 `internal/domain/myservice/service/myservice.go`：

```go
package service

import (
	"context"
	"fmt"
)

// Service 我的服务
type Service struct {
	// 依赖项
}

// NewService 创建服务实例
func NewService() *Service {
	return &Service{}
}

// DoSomething 执行某个操作
func (s *Service) DoSomething(ctx context.Context, param string) (string, error) {
	// 业务逻辑
	if param == "" {
		return "", fmt.Errorf("param is required")
	}

	result := fmt.Sprintf("Processed: %s", param)
	return result, nil
}
```

#### 步骤 5: 实现 Handler 层

创建 `internal/domain/myservice/myservice_handler.go`：

```go
package myservice

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	myservicev1 "github.com/HJH0924/go-template-project/sdk/go/myservice/v1"
	"github.com/HJH0924/go-template-project/internal/domain/myservice/service"
)

// Handler 处理器
type Handler struct {
	service *service.Service
	logger  *slog.Logger
}

// NewHandler 创建处理器
func NewHandler(svc *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: svc,
		logger:  logger,
	}
}

// DoSomething 处理请求
func (h *Handler) DoSomething(
	ctx context.Context,
	req *connect.Request[myservicev1.DoSomethingRequest],
) (*connect.Response[myservicev1.DoSomethingResponse], error) {
	h.logger.InfoContext(ctx, "processing request",
		slog.String("param", req.Msg.GetParam()))

	result, err := h.service.DoSomething(ctx, req.Msg.GetParam())
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to process",
			slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&myservicev1.DoSomethingResponse{
		Result: result,
	}), nil
}
```

#### 步骤 6: 在 serve.go 中注册

编辑 `internal/cmd/serve.go`，添加服务注册：

```go
import (
	myservicev1connect "github.com/HJH0924/go-template-project/sdk/go/myservice/v1/myservicev1connect"
	"github.com/HJH0924/go-template-project/internal/domain/myservice"
	myservice_service "github.com/HJH0924/go-template-project/internal/domain/myservice/service"
)

// 在 RunE 函数中
myService := myservice_service.NewService()
myHandler := myservice.NewHandler(myService, logger)
path, handler := myservicev1connect.NewMyServiceHandler(myHandler)
mux.Handle(path, handler)
```

#### 步骤 7: 编写单元测试

创建 `internal/domain/myservice/service/myservice_test.go`：

```go
package service

import (
	"context"
	"testing"
)

func TestService_DoSomething(t *testing.T) {
	tests := []struct {
		name    string
		param   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid param",
			param:   "test",
			want:    "Processed: test",
			wantErr: false,
		},
		{
			name:    "empty param",
			param:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.DoSomething(context.Background(), tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSomething() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DoSomething() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

### 2. 工作流程

推荐的开发工作流：

```bash
# 1. 创建功能分支
git checkout -b feature/my-new-feature

# 2. 定义 proto 接口
vim proto/myservice/v1/myservice.proto

# 3. 生成代码
make generate

# 4. 实现 Service 和 Handler
# ... 编写代码 ...

# 5. 编写测试
vim internal/domain/myservice/service/myservice_test.go

# 6. 运行测试
make test

# 7. 开发模式运行验证
make dev

# 8. 手动测试 API
curl --request POST \
  --url http://localhost:8080/myservice.v1.MyService/DoSomething \
  --header 'Content-Type: application/json' \
  --data '{"param": "test"}'

# 9. 提交代码（会自动触发 pre-commit hook）
git add .
git commit -m "feat: add MyService"

# 10. 推送到远程
git push origin feature/my-new-feature
```

## API 开发

### Connect RPC

项目使用 [Connect RPC](https://connectrpc.com/) 框架，它提供：

- **协议兼容**: 支持 gRPC、gRPC-Web 和自定义的 Connect 协议
- **HTTP/1.1 支持**: 无需 HTTP/2
- **标准 HTTP**: 可以使用任何 HTTP 客户端
- **流式支持**: 支持服务器流、客户端流和双向流
- **浏览器友好**: 可以直接从浏览器调用

### 请求示例

**使用 curl:**

```bash
curl --request POST \
  --url http://localhost:8080/user.v1.UserService/SayHello \
  --header 'Content-Type: application/json' \
  --data '{"name": "Alice"}'
```

**使用 Go 客户端:**

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	userv1 "github.com/HJH0924/go-template-project/sdk/go/user/v1"
	"github.com/HJH0924/go-template-project/sdk/go/user/v1/userv1connect"
)

func main() {
	client := userv1connect.NewUserServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	req := connect.NewRequest(&userv1.SayHelloRequest{
		Name: "Alice",
	})

	resp, err := client.SayHello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Msg.GetMessage())
}
```

### 错误处理

Connect 提供了标准的错误码：

```go
// 参数错误
return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid param"))

// 未找到
return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("not found"))

// 内部错误
return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("internal error"))

// 未授权
return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("unauthorized"))
```

## 测试

### 单元测试

```bash
# 运行所有测试（包含竞态检测和覆盖率）
make test

# 快速测试（不包含竞态检测）
make test-short

# 运行特定包的测试
go test -v ./internal/domain/user/service/...

# 运行特定测试
go test -v -run TestService_SayHello ./internal/domain/user/service/
```

### 测试覆盖率

```bash
# 生成覆盖率报告
make test

# 查看覆盖率 HTML 报告
open coverage.html
```

### 编写测试

**Service 层测试示例:**

```go
package service

import (
	"context"
	"testing"
)

func TestService_SayHello(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "with name",
			arg:  "Alice",
			want: "Hello, Alice!",
		},
		{
			name: "empty name",
			arg:  "",
			want: "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got := s.SayHello(context.Background(), tt.arg)
			if got != tt.want {
				t.Errorf("SayHello() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

### 集成测试

```bash
# 启动服务
make dev &

# 等待服务启动
sleep 2

# 测试健康检查
curl http://localhost:8080/health

# 测试 API
curl --request POST \
  --url http://localhost:8080/user.v1.UserService/SayHello \
  --header 'Content-Type: application/json' \
  --data '{"name": "Test"}'

# 停止服务
pkill -f "go run ./cmd/server"
```

## 配置管理

项目使用 [Viper](https://github.com/spf13/viper) 进行配置管理。

### 配置文件

配置文件位于 `configs/config.yaml`：

```yaml
server:
  host: "0.0.0.0"  # 监听地址
  port: 8080        # 监听端口
```

### 加载配置

配置在 `internal/config/config.go` 中定义和加载：

```go
package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	instance *Config
	once     sync.Once
)

// Config 全局配置
type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Load 加载配置文件
func Load(configPath string) error {
	var err error
	once.Do(func() {
		v := viper.New()
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")

		if err = v.ReadInConfig(); err != nil {
			err = fmt.Errorf("failed to read config: %w", err)
			return
		}

		instance = &Config{}
		if err = v.Unmarshal(instance); err != nil {
			err = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}
	})
	return err
}

// Get 获取配置实例
func Get() *Config {
	if instance == nil {
		panic("config not loaded")
	}
	return instance
}
```

### 使用配置

```go
import "github.com/HJH0924/go-template-project/internal/config"

// 加载配置
if err := config.Load(configPath); err != nil {
	log.Fatal(err)
}

// 使用配置
cfg := config.Get()
addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
```

### 环境变量

可以通过环境变量覆盖配置：

```bash
# 设置环境变量
export SERVER_HOST=127.0.0.1
export SERVER_PORT=9090

# 运行服务
./bin/server serve --config configs/config.yaml
```

## 日志记录

项目使用 Go 1.21+ 的 `log/slog` 包进行结构化日志记录。

### 日志级别

- **Debug**: 详细的调试信息
- **Info**: 一般信息
- **Warn**: 警告信息
- **Error**: 错误信息

### 使用示例

```go
import "log/slog"

// Info 级别
logger.InfoContext(ctx, "processing request",
	slog.String("name", name),
	slog.Int("id", 123))

// Error 级别
logger.ErrorContext(ctx, "failed to process request",
	slog.Any("error", err),
	slog.String("detail", "additional info"))

// Debug 级别
logger.DebugContext(ctx, "debug information",
	slog.Any("data", data))

// Warn 级别
logger.WarnContext(ctx, "warning message",
	slog.String("reason", "something unusual"))
```

### 日志配置

日志配置在 `internal/cmd/serve.go` 中：

```go
// 创建 JSON 格式的结构化日志
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.String("timestamp", a.Value.Time().Format(time.RFC3339))
		}
		return a
	},
}))
```

### 日志输出示例

```json
{
  "timestamp": "2025-01-29T10:30:45Z",
  "level": "INFO",
  "msg": "processing request",
  "name": "Alice",
  "id": 123
}
```

## 故障排除

### 开发工具缺失

**问题**：
```
command not found: golangci-lint
command not found: gofumpt
```

**解决方法**：
```bash
# 自动安装所有开发依赖
make install-deps

# 验证安装
golangci-lint --version
gofumpt -version
buf --version
```

### Proto 生成问题

**问题**：`make generate` 失败

**解决方法**：
```bash
# 检查 buf 安装
buf --version

# 检查 buf.yaml 配置
buf lint

# 清理并重新生成
rm -rf sdk/go
make generate
```

### 构建错误

**问题**：编译失败或依赖问题

**解决方法**：
```bash
# 清理缓存
go clean -cache

# 整理依赖
go mod tidy

# 更新依赖
go get -u ./...

# 重新构建
make build
```

### Pre-commit Hook 问题

**问题**：Hook 未生效或提交被阻止

**解决方法**：
```bash
# 重新安装 hook
make install-hooks

# 验证 hook 文件存在
ls -la .git/hooks/pre-commit

# 手动运行检查
make format
make lint

# 如需临时跳过 hook（不推荐）
git commit --no-verify -m "message"
```

### 测试失败

**问题**：测试运行失败

**解决方法**：
```bash
# 运行测试并显示详细输出
go test -v ./...

# 运行特定测试
go test -v -run TestName ./internal/domain/user/service/

# 清理测试缓存
go clean -testcache
go test ./...

# 检查竞态条件
go test -race ./...
```

### 服务无法启动

**问题**：服务启动失败

**解决方法**：
```bash
# 检查端口是否被占用
lsof -i :8080

# 检查配置文件是否正确
cat configs/config.yaml

# 使用调试模式运行
go run -race ./cmd/server serve --config configs/config.yaml

# 检查日志输出
tail -f /tmp/server.log
```

### Import 路径错误

**问题**：`cannot find package` 或 import 路径错误

**解决方法**：
```bash
# 确认模块名称
go mod edit -print | grep module

# 更新所有 import 路径
find . -type f -name "*.go" -exec sed -i '' 's|old-module|new-module|g' {} +

# 重新整理依赖
go mod tidy

# 运行格式化
make format
```

## Makefile 命令参考

```bash
make help              # 显示所有可用命令
make install-deps      # 安装开发依赖和 pre-commit hook
make install-hooks     # 仅安装 Git pre-commit hooks
make format            # 格式化代码
make lint              # 运行代码检查
make build             # 构建项目（包含 format 和 lint）
make run               # 构建并运行服务器
make dev               # 开发模式运行（无需构建）
make test              # 运行测试并生成覆盖率报告
make test-short        # 快速测试（不包含竞态检测）
make clean             # 清理构建产物
make generate          # 生成 proto 代码
make docs              # 运行文档服务器
```

## 最佳实践

### 1. 代码组织

- 使用领域驱动设计（DDD）组织代码
- Service 层包含业务逻辑，Handler 层处理 RPC 请求
- 每个域独立，降低耦合

### 2. 错误处理

- 使用 `fmt.Errorf` 包装错误，保留错误链
- 在 Service 层返回业务错误
- 在 Handler 层转换为 Connect 错误码

### 3. 测试

- 编写单元测试，覆盖核心业务逻辑
- 使用表驱动测试（table-driven tests）
- 保持测试独立，避免依赖外部资源

### 4. 日志

- 使用结构化日志，便于查询和分析
- 记录关键操作和错误信息
- 避免记录敏感信息（密码、密钥等）

### 5. 配置

- 使用配置文件管理环境相关配置
- 支持环境变量覆盖
- 不在代码中硬编码配置值

### 6. 依赖管理

- 定期更新依赖到最新稳定版本
- 使用 `go mod tidy` 清理未使用的依赖
- 提交 `go.sum` 到版本控制

## 贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
   - Pre-commit hook 会自动运行格式化和检查
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat:` 新功能
- `fix:` 修复 bug
- `docs:` 文档更新
- `style:` 代码格式（不影响代码运行）
- `refactor:` 重构代码
- `test:` 测试相关
- `chore:` 构建过程或辅助工具的变动

## 资源

- [Connect RPC 文档](https://connectrpc.com/docs/)
- [Protocol Buffers 指南](https://protobuf.dev/)
- [Go slog 包](https://pkg.go.dev/log/slog)
- [Buf CLI 文档](https://buf.build/docs/)
- [Viper 文档](https://github.com/spf13/viper)
- [Cobra 文档](https://github.com/spf13/cobra)
- [golangci-lint 文档](https://golangci-lint.run/)
