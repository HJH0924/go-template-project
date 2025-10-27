# go-template-project

一个现代化的 Go 项目模板，展示 Go 语言开发的最佳实践。

## 特性

- 🚀 **gRPC/Connect** - 使用 Connect RPC 框架，支持 HTTP/1.1、HTTP/2 和 gRPC
- 📝 **结构化日志** - 使用 Go 1.21+ 的 `log/slog` 包进行结构化日志记录
- ⚙️ **配置管理** - 使用 Viper 读取 YAML 配置文件，支持全局单例模式
- 🎯 **命令行工具** - 使用 Cobra 构建强大的 CLI 工具
- 🏗️ **分层架构** - Handler、Service 层清晰分离，易于维护和测试
- 🧪 **单元测试** - 包含完整的单元测试示例
- 🔍 **代码检查** - 集成 golangci-lint，确保代码质量
- 📦 **Protocol Buffers** - 使用 Buf 管理 proto 文件和代码生成

## 快速开始

### 前置要求

- Go 1.24.7 或更高版本
- Make

### 安装开发依赖

```bash
make install-deps
```

这将安装以下工具：
- golangci-lint - 代码检查工具
- gofumpt - 代码格式化工具
- goimports - import 排序工具
- goimports-reviser - import 分组工具
- buf - Protocol Buffers 工具

### 构建项目

```bash
make build
```

### 运行服务器

```bash
# 方式 1: 构建后运行
make run

# 方式 2: 开发模式（直接运行）
make dev

# 方式 3: 直接运行二进制文件
./bin/server serve --config configs/config.yaml
```

服务器将在 `0.0.0.0:8080` 启动。

### 运行测试

```bash
# 运行所有测试（包含竞态检测和覆盖率）
make test

# 快速测试（不包含竞态检测）
make test-short
```

### 生成 Proto 代码

```bash
make generate
```

## 项目结构

```
.
├── cmd/
│   └── server/           # 服务器入口
│       └── main.go
├── configs/              # 配置文件
│   └── config.yaml
├── internal/
│   ├── cmd/             # Cobra 命令定义
│   │   ├── root.go      # 根命令
│   │   └── serve.go     # serve 子命令
│   ├── config/          # 配置管理
│   │   └── config.go    # Viper 配置加载
│   └── domain/          # 业务领域
│       └── user/
│           ├── service/        # 业务逻辑层
│           │   ├── user.go     # Service 接口和实现
│           │   └── user_test.go # 单元测试
│           └── user_handler.go  # Handler 层（Connect RPC）
├── proto/               # Protocol Buffers 定义
│   └── user/
│       └── v1/
│           └── user.proto
├── sdk/                 # 生成的 SDK 代码
│   └── go/
│       └── user/
│           └── v1/
├── scripts/             # 脚本文件
│   └── pre-commit       # Git pre-commit hook
├── buf.gen.yaml         # Buf 代码生成配置
├── buf.yaml             # Buf 配置
├── go.mod               # Go 模块定义
├── Makefile             # Make 命令定义
└── README.md
```

## API 文档

### SayHello

向指定的人打招呼。

**请求:**
```bash
curl --request POST \
  --url http://localhost:8080/user.v1.UserService/SayHello \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "Alice"
  }'
```

**响应:**
```json
{
  "message": "Hello, Alice!"
}
```

如果不提供 `name` 参数，默认返回 `Hello, World!`。

### 健康检查

```bash
curl http://localhost:8080/health
```

**响应:**
```
OK
```

## 开发指南

### 添加新的 Service

1. 在 `proto/` 目录下定义 proto 文件
2. 运行 `make generate` 生成代码
3. 在 `internal/domain/` 下创建对应的目录
4. 实现 Service 层业务逻辑
5. 实现 Handler 层处理 RPC 请求
6. 在 `internal/cmd/serve.go` 中注册服务

### 代码规范

项目使用 golangci-lint 进行代码检查，配置文件为 `.golangci.toml`。

运行代码检查：
```bash
make lint
```

格式化代码：
```bash
make format
```

### 提交代码

项目配置了 pre-commit hook，在提交代码前会自动进行：
- 代码格式化
- 代码检查
- 单元测试

如果检查不通过，提交将被阻止。

## 技术栈

- **Web 框架**: [Connect](https://connectrpc.com/) - 现代化的 RPC 框架
- **配置管理**: [Viper](https://github.com/spf13/viper) - 配置解决方案
- **命令行**: [Cobra](https://github.com/spf13/cobra) - CLI 框架
- **日志**: [slog](https://pkg.go.dev/log/slog) - Go 标准库结构化日志
- **Proto 工具**: [Buf](https://buf.build/) - Protocol Buffers 工具链
- **代码检查**: [golangci-lint](https://golangci-lint.run/) - Go linters 聚合器

## Make 命令

```bash
make help              # 显示所有可用命令
make install-deps      # 安装开发依赖
make install-hooks     # 安装 git hooks
make format            # 格式化代码
make lint              # 运行代码检查
make build             # 构建项目
make run               # 运行服务器
make dev               # 开发模式运行
make test              # 运行测试
make test-short        # 快速测试
make clean             # 清理构建产物
make generate          # 生成 proto 代码
```

## 配置说明

配置文件位于 `configs/config.yaml`：

```yaml
server:
  host: "0.0.0.0"  # 监听地址
  port: 8080        # 监听端口
```

可以通过 `--config` 参数指定不同的配置文件：

```bash
./bin/server serve --config /path/to/config.yaml
```

## License

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！
